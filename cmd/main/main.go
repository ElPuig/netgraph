package main

import (
	"net/http"
	"text/template"
)

var netview_path string = "../../internal/netview/"

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	index_template, err := template.ParseFiles(netview_path + "index.html")
	if err != nil {
		panic(err)
	} else {
		index_template.Execute(rw, nil)
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe("localhost:3000", nil)
}
