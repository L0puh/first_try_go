package main

import (
	"html/template"
	"net/http"
)
type Page struct {
    Title string 
    Info string 
}
func index(w http.ResponseWriter, r *http.Request ) {
    pg := Page{Title: "some title", Info: "some random info"}
    tmpl, _ := template.ParseFiles("index.html")
    tmpl.Execute(w, pg)
}
func main () {
    http.HandleFunc("/", index)
    http.ListenAndServe(":9000", nil)
}
