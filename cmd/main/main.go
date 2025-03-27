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
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/ElPuig/netgraph/pkg/graph_vis"
	"github.com/ElPuig/netgraph/pkg/xml_loader"
	"github.com/alexflint/go-arg"
)

var node_map graph_vis.NodeMap

var args struct {
	Source string `arg:"positional,required"`
}

var netview_path string = "internal/netview/"

type PageData struct {
	Date  string
	Nodes json.RawMessage
	Edges string
}

var pageData PageData

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	index_template, err := template.ParseFiles(netview_path + "index.html")
	if err != nil {
		panic(err)
	}
	index_template.Execute(rw, pageData)
}

func main() {
	arg.MustParse(&args)
	// TODO: check Source arg and add '/' at the end to avoid errors later on.

	xml_files, err := xml_loader.DownloadXmlFiles(args.Source, `\d+\.\d+\.\d+\.\d+`) // TODO: regex hardcoded. Turn into argument in the future?
	if err != nil {
		panic(err)
	}

	node_map = graph_vis.GetNodeMap(xml_files)

	pageData.Date = time.Now().String()
	pageData.Nodes = node_map.ToVisJson()
	pageData.Edges = ""

	css_fs := http.FileServer(http.Dir(netview_path + "css"))
	http.Handle("/css/", http.StripPrefix("/css/", css_fs))
	http.HandleFunc("/", indexHandler)
	fmt.Println("Listening on localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}
