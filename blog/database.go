package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db* sql.DB
type Post struct {
    Id int
    Title string
    Content string
}

func connect_db () {
    conn, _ := godotenv.Read(".env")
    var err error
    db, err = sql.Open("postgres", conn["conn"])
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("connected to db")
}
func add_post(title string, content string) {
    _, err := db.Exec(`insert into blog (title, content) values($1, $2)`, title, content)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("new post: %s", title)

}

func get_posts(limit int) []Post {
    rows, err := db.Query(`SELECT * FROM blog LIMIT $1`, limit) 
    defer rows.Close()
    var posts[]Post
    var(id int; title string; content string)
    if err != nil {
        log.Fatal(err)
    }
    for rows.Next() {
        rows.Scan(&id, &title, &content)
        posts = append(posts, Post{id, title, content}) 
    }
    return posts
}
func close_db() {
    db.Close()
    log.Printf("db is closed")
}
func delete_postById(id int ) {
    _, err := db.Exec(`DELETE FROM blog WHERE id = $1`, id)
    if err != nil {
        log.Fatal(err)
    }
}

func get_post(id int) Post {
    rows, err := db.Query(`SELECT * FROM blog WHERE id=$1`, id) 
    defer rows.Close()

    var(id_p int; title string; content string)
    if err != nil {
        log.Fatal(err)
    }

    for rows.Next() {
        rows.Scan(&id_p, &title, &content)
    }
    post := Post{id_p, title, content}
    return post
}

