package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var tmp *template.Template

func main () {
    connect_db()
    style := http.FileServer(http.Dir("styles"))
    http.Handle("/styles/", http.StripPrefix("/styles/", style))
    
    var err error
    tmp, err = template.ParseGlob("tmp/*.html")
    print_err(err)

    http.HandleFunc("/", home) 
    http.HandleFunc("/add/", add) 
    http.HandleFunc("/delete/", delete_post) 
    http.HandleFunc("/post/", show_post) 
    http.HandleFunc("/signup/", sign_up) 
    http.HandleFunc("/login/", log_in) 
    http.HandleFunc("/profile/", profile) 
    http.HandleFunc("/logout/", logout)
    log.Fatal(http.ListenAndServe("127.0.0.1:9000", nil))

    defer close_db()
}

func print_err(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func profile(w http.ResponseWriter, r *http.Request) {
    name := get_cookie_username(r)
    tmp.ExecuteTemplate(w, "profile.html", name)
}

func sign_up (w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        username := r.FormValue("username")
        email    := r.FormValue("email")
        password := r.FormValue("password")
        pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)

        create_cookie(w, "login", username)
        print_err(err)
        if create_user(username, email, pass) != false {
            http.Redirect(w, r, "/login/", http.StatusFound)
        }
    }
    tmp.ExecuteTemplate(w, "sign_up.html", nil)
}
func is_cookie(r *http.Request) bool {
    _, err := r.Cookie("login")
    if err != nil {
        return false 
    }
    return true
}
func get_cookie_username(r *http.Request) string  {
    if is_cookie(r) {
        cookie, _ := r.Cookie("login") 
        return cookie.Value
    }
    return ""
}
func create_cookie(w http.ResponseWriter, name string, value string) {
    cookie := http.Cookie {
        Name: name,
        Value: value,
        Path: "/",
    } 
    http.SetCookie(w, &cookie)
}

func logout(w http.ResponseWriter, r *http.Request) {
    if is_cookie(r) {
        c := http.Cookie {
            Name:       "login", 
            Value:      "",
            Path:       "/",
            Expires:    time.Unix(0, 0),
        }
        http.SetCookie(w, &c)
    }
    http.Redirect(w, r, "/profile", http.StatusFound)
}

func log_in (w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        username := r.FormValue("username")
        email    := r.FormValue("email")
        password := r.FormValue("password")

        if (check_exist(username, email)) {
            hash_pass := get_password(username)
            if (bcrypt.CompareHashAndPassword(hash_pass, []byte(password)) == nil) {
                create_cookie(w, "login", username)
                http.Redirect(w, r, "/", http.StatusFound)
            } 
        }
    }
    tmp.ExecuteTemplate(w, "log_in.html", nil)
}


func home(w http.ResponseWriter, r *http.Request) {
    posts := get_posts(10)
    tmp.ExecuteTemplate(w, "home.html", posts)
}

func add(w http.ResponseWriter, r *http.Request) {
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
    print_err(err)
    return num
}

func delete_post(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/delete/"):]

    delete_postById(convert(id))
    log.Printf("deleted: post %s", id)
    http.Redirect(w, r, "/", http.StatusFound)
}

func show_post(w http.ResponseWriter, r *http.Request) {
    id   := r.URL.Path[len("/post/"):]
    post := get_post(convert(id))
    id_post := convert(id) 
    if r.Method == "POST" {
        if is_cookie(r) {
            name := get_cookie_username(r)
            id_user := get_id(name)
            comment := r.FormValue("comment") 
            create_comment(id_user, id_post, comment)
        } else {
            log.Printf("cookie can't be found")
        }
    }
    comments := get_comments(id_post)
    post.Comments = comments

    tmp.ExecuteTemplate(w, "post.html", post)
}

