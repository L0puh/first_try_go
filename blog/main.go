package main

import (
	"html/template"
	"log"
	"net/http"
)

var tmp *template.Template
func main () {
    connect_db()
    style := http.FileServer(http.Dir("styles"))
    http.Handle("/styles/", http.StripPrefix("/styles/", style))

    tmp, _ = template.ParseGlob("tmp/*.html")
    http.HandleFunc("/", home) 
    http.HandleFunc("/add", add) 
    log.Fatal(http.ListenAndServe("127.0.0.1:9000", nil))
}

func home( w http.ResponseWriter, r *http.Request ) {
    tmp.ExecuteTemplate(w, "home.html", nil)
}

func add(w http.ResponseWriter, r *http.Request ) {
    //TODO
}
