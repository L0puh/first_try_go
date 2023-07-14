package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

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
    http.ListenAndServe("127.0.0.1:9000", nil)
}

func posts(w http.ResponseWriter, r *http.Request){
    ins := `select * from posts`
    rows, _ := db.Query(ins)
    
    var (title, info string; id int)
    for rows.Next() {
        _ = rows.Scan(&id, &title, &info)
        posts := Post{id, title, info}
        tmp.ExecuteTemplate(w, "posts.html", posts)
    }

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

