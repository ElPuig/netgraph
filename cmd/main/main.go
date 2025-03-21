/*
Copyright 2025 David García De Mercado.

This file is part of Netgraph.

   Netgraph is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   Netgraph is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with Netgraph.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexflint/go-arg"
)

var args struct {
	Source string `arg:"positional,required"`
}

type RequestXMLData struct {
	XMLName xml.Name `xml:"REQUEST"`
	Device  struct {
		Info struct {
			Name     string `xml:"NAME"`
			Model    string `xml:"MODEL"`
			Location string `xml:"LOCATION"`
		} `xml:"INFO"`
		Ports []struct {
			IfName      string `xml:"IFNAME"`
			Connections struct {
				Cdp        string `xml:"CDP"`
				Connection []struct {
					SysName string   `xml:"SYSNAME"`
					Macs    []string `xml:"MAC"`
				} `xml:"CONNECTION"`
			} `xml:"CONNECTIONS"`
			Vlans []struct {
				Number string `xml:"NUMBER"`
				Name   string `xml:"NAME"`
				Tagged string `xml:"TAGGED"`
			} `xml:"VLANS>VLAN"`
		} `xml:"PORTS>PORT"`
	} `xml:"CONTENT>DEVICE"`
}

func getUrlBody(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Could not access '%s': %v", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Could not access '%s': %d", url, resp.StatusCode)
	}

	return resp.Body, nil
}

func getFileLinks(url string, regex string) ([]string, error) {
	body, err := getUrlBody(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	r, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	var res []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			if r.MatchString(href) {
				res = append(res, url+href)
			}
		}
	})
	return res, nil
}

func getXmlFromUrl(url string) (RequestXMLData, error) {
	body, err := getUrlBody(url)
	if err != nil {
		return RequestXMLData{}, err
	}
	defer body.Close()

	content, err := io.ReadAll(body)
	if err != nil {
		return RequestXMLData{}, fmt.Errorf("Error reading XML file: %v", err)
	}

	var xml_data RequestXMLData
	err = xml.NewDecoder(bytes.NewReader(content)).Decode(&xml_data)
	if err != nil {
		return RequestXMLData{}, fmt.Errorf("Error reading XML file: %v", err)
	}
	return xml_data, nil
}

func downloadXmlFiles(url string, regex string) ([]RequestXMLData, error) {
	links, err := getFileLinks(url, regex)
	if err != nil {
		return nil, err
	}

	res := []RequestXMLData{}
	for _, lnk := range links {
		xml_data, err := getXmlFromUrl(lnk)
		if err != nil {
			return nil, err
		}
		res = append(res, xml_data)
	}
	return res, nil // TODO: Return an actual value. Remove print above.
}

var netview_path string = "internal/netview/"

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	index_template, err := template.ParseFiles(netview_path + "index.html")
	if err != nil {
		panic(err)
	}
	index_template.Execute(rw, nil)
}

func main() {
	arg.MustParse(&args)
	// TODO: check Source arg and add '/' at the end to avoid errors later on.

	xml_files, err := downloadXmlFiles(args.Source, `\d+\.\d+\.\d+\.\d+`) // TODO: regex hardcoded. Turn into argument in the future?
	if err != nil {
		panic(err)
	}
	fmt.Println(xml_files)

	css_fs := http.FileServer(http.Dir(netview_path + "css"))
	http.Handle("/css/", http.StripPrefix("/css/", css_fs))
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe("localhost:3000", nil)
}
