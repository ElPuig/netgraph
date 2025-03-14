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
	"html/template"
	"net/http"

	"github.com/alexflint/go-arg"
)

var args struct {
	Source string `arg:"positional,required"`
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

	css_fs := http.FileServer(http.Dir(netview_path + "css"))
	http.Handle("/css/", http.StripPrefix("/css/", css_fs))
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe("localhost:3000", nil)
}
