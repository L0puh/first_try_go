package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db* sql.DB;

func connect_db () {
    conn, _ := godotenv.Read(".env")
    db, _ = sql.Open("postgress", conn["conn"])
    fmt.Println("+ connected")
}

