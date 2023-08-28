package main

import (
	"database/sql"
	"fmt"
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
    ins := fmt.Sprintf(`insert into blog (title, content)
                        values('%s', '%s')`,  title, content)

    _, err := db.Exec(ins)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("new post: %s", title)

}

func get_posts(limit int) []Post {
    ins := fmt.Sprintf(`SELECT * FROM blog LIMIT %d`, limit)
    rows, err := db.Query(ins) 
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
    ins := fmt.Sprintf(`DELETE FROM blog WHERE id = %d`, id)
    _, err := db.Exec(ins)
    if err != nil {
        log.Fatal(err)
    }
}
