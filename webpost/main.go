package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Page struct {
    Title string 
    Info string 
    Status int 
    List[]string
}
type Post struct {
    Id int 
    Title string 
    Info string 
}
func connect_db() {
    var envMap map[string]string;
    envMap,_ = godotenv.Read(".env")

    psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", envMap["host"], 5432, envMap["user"], envMap["password"], envMap["dbname"])
    db, _ = sql.Open("postgres", psqlconn)
    fmt.Println("connected")
}
var db *sql.DB
var tmp *template.Template

func main () {
    connect_db()
    fs := http.FileServer(http.Dir("styles"))
    http.Handle("/styles/", http.StripPrefix("/styles/", fs)) // serve css styles 
    tmp, _ = template.ParseGlob("templates/*.html") // parse templates from the dir
    http.HandleFunc("/add/", add)
    http.HandleFunc("/posts/", posts)
    http.HandleFunc("/delete/", delete_post)
    http.ListenAndServe("127.0.0.1:9000", nil)
}

func get_all_posts() []Post{
    ins := `select * from posts`
    rows, _ := db.Query(ins)
    var posts_list[]Post 
    var (title, info string; id int)
    for rows.Next() {
        _ = rows.Scan(&id, &title, &info)
        post := Post{id, title, info}
        posts_list = append(posts_list, post)
    }
    return posts_list
}
func posts(w http.ResponseWriter, r *http.Request){
    posts := get_all_posts()
    // for post := 0; post < len(posts); post++ {
    //     tmp.ExecuteTemplate(w, "posts.html", posts)
    // }

    tmp.ExecuteTemplate(w, "posts.html", posts)
}
func add(w http.ResponseWriter, r *http.Request ) {
    if r.Method == "POST" {
        title := r.FormValue("title")
        text := r.FormValue("text")
        ins := fmt.Sprintf(`insert into posts(title, info) values('%s', '%s')`, title, text)
        _, err := db.Exec(ins)
        if err != nil{
            panic(err)
        }
        http.Redirect(w, r, "/posts/", http.StatusFound)
    }
    tmp.ExecuteTemplate(w, "add.html", nil)
}
func delete_post(w http.ResponseWriter, r *http.Request){
    if r.Method == "POST" {
        id := r.FormValue("id_post")
        id_post, _ := strconv.Atoi(id)
        fmt.Printf("%d is deleted", id_post)
        ins := fmt.Sprintf(`delete from posts WHERE id=%d`, id_post)
        _, err := db.Exec(ins)
        if err != nil {
            panic(err)
        }
        http.Redirect(w, r, "/posts/", http.StatusFound)
    }
    posts := get_all_posts()
    for i:=0; i < len(posts); i++ {
        tmp.ExecuteTemplate(w, "delete.html", posts[i])
    }
}
