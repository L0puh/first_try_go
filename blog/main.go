package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmp *template.Template
func main () {
    connect_db()
    style := http.FileServer(http.Dir("styles"))
    http.Handle("/styles/", http.StripPrefix("/styles/", style))

    tmp, _ = template.ParseGlob("tmp/*.html")
    http.HandleFunc("/", home) 
    http.HandleFunc("/add/", add) 
    http.HandleFunc("/delete/", delete_post) 
    http.HandleFunc("/post/", show_post) 
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
        http.Redirect(w, r, "/", http.StatusFound)
    }
    tmp.ExecuteTemplate(w, "post_add.html", nil)
}

func convert(id string) int {
    num, err := strconv.Atoi(id)
    if err != nil {
        log.Fatal(err)
    }
    return num
}
func delete_post(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/delete/"):]

    delete_postById(convert(id))
    log.Printf("deleted: post %s", id)
    http.Redirect(w, r, "/", http.StatusFound)
}

func show_post(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/post/"):]
    post := get_post(convert(id))
    tmp.ExecuteTemplate(w, "post.html", post)
}

