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
	"fmt"
	"html/template"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexflint/go-arg"
)

var args struct {
	Source string `arg:"positional,required"`
}

func getFileLinks(url string, regex string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Could not access '%s': %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Could not access '%s': %d", url, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
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

func downloadXmlFiles(url string, regex string) (string, error) {
	links, err := getFileLinks(url, regex)
	if err != nil {
		return "", err
	}
	fmt.Println(links)
	return "", nil // TODO: Return an actual value. Remove print above.
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
