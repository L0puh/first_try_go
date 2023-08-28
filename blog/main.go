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
    http.HandleFunc("/add/", add) 
    log.Fatal(http.ListenAndServe("127.0.0.1:9000", nil))
    close_db()
}

func home( w http.ResponseWriter, r *http.Request ) {
    posts := get_posts(10)
    tmp.ExecuteTemplate(w, "home.html", posts)
}

func add(w http.ResponseWriter, r *http.Request ) {
    if r.Method == "POST" {
        title := r.FormValue("title")
        content := r.FormValue("content")
        add_post(title, content)
        // log.Printf("new post: %s", title)
        http.Redirect(w, r, "/", http.StatusFound)
    }
    tmp.ExecuteTemplate(w, "post_add.html", nil)
}
