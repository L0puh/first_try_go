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
    Comments []Comment
}
type Comment struct {
    User string
    Comment string
}
func connect_db () {
    conn, _ := godotenv.Read(".env")
    var err error
    db, err = sql.Open("postgres", conn["conn"])
    print_err(err)
    log.Printf("connected to db")
}
func add_post(title string, content string) {
    _, err := db.Exec(`insert into blog (title, content) values($1, $2)`, title, content)
    print_err(err)
    log.Printf("new post: %s", title)

}

func create_user(username string, email string, password []byte) bool {
    if check_exist(username, email) {
        return false
    }
    _, err := db.Exec(`insert into users(username, 
                email, password) values($1, $2, $3)`, username, email, password)
    print_err(err)
    return true
}

func check_exist(username string, email string) bool {
    var res bool 
    err := db.QueryRow(`SELECT EXISTS (SELECT * FROM users 
            WHERE username = $1 OR email = $2)`, username, email).Scan(&res)
    print_err(err)
    return res
}
func get_password(username string) []byte{
    rows, err := db.Query(`SELECT password FROM users WHERE username=$1`, username) 
    defer rows.Close()
    var pass []byte
    print_err(err)
    for rows.Next() {
        rows.Scan(&pass)
    }
    return pass 
}
func get_posts(limit int) []Post {
    rows, err := db.Query(`SELECT * FROM blog LIMIT $1`, limit) 
    defer rows.Close()
    var posts[]Post
    var(id int; title string; content string)
    print_err(err)
    for rows.Next() {
        rows.Scan(&id, &title, &content)
        posts = append(posts, Post{id, title, content, []Comment{}}) 
    }
    return posts
}
func close_db() {
    db.Close()
    log.Printf("db is closed")
}
func delete_postById(id int ) {
    _, err := db.Exec(`DELETE FROM blog WHERE id = $1`, id)
    print_err(err)
}


func get_id(name string) int {
    rows, err := db.Query(`SELECT id FROM users WHERE username=$1`, name) 
    defer rows.Close()
    var id int 
    print_err(err)
    for rows.Next() {
        rows.Scan(&id)
    }
    return id
}

func create_comment(user_id int, post_id int, comment string) {
    _, err := db.Exec(`insert into comments(user_id, post_id, comment) values($1, $2, $3)`,
    user_id, post_id, comment)
    print_err(err)
    log.Printf("new comment: %s", comment)
}

func get_comments(post_id int) []Comment{
    rows, err := db.Query(`SELECT user_id, comment FROM comments WHERE post_id=$1`, post_id) 
    print_err(err)
    
    defer rows.Close()
    var comments[]Comment 
    var(user int; comment string)

    for rows.Next() {
        rows.Scan(&user, &comment)
        comments = append(comments, Comment{get_name(user), comment})
    }
    return comments

}
func get_name(id int) string {
    rows, err := db.Query(`SELECT username FROM users WHERE id=$1`, id) 
    defer rows.Close()
    var name string 
    print_err(err)
    for rows.Next() {
        rows.Scan(&name)
    }
    return name
}
func get_post(id int) Post {
    rows, err := db.Query(`SELECT * FROM blog WHERE id=$1`, id) 
    defer rows.Close()

    var(id_p int; title string; content string)
    print_err(err)

    for rows.Next() {
        rows.Scan(&id_p, &title, &content)
    }

    post := Post{id_p, title, content, []Comment{}}
    return post
}

