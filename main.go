package main

import (
	"html/template"
	"net/http"
)
type Page struct {
    Title string 
    Info string 
}
func main () {
    fs := http.FileServer(http.Dir("styles"))
    http.Handle("/styles/", http.StripPrefix("/styles/", fs))

    http.HandleFunc("/", index)
    http.ListenAndServe(":9000", nil)
}

func index(w http.ResponseWriter, r *http.Request ) {
    infos := []Page{
        {"first step", "info for first one"},
        {"second step", "info for second one"},
    }
    tmpl, _ := template.ParseFiles("templates/index.html")
    for i := range infos {
        tmpl.Execute(w, infos[i])
    }
}
