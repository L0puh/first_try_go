package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request ) {
    fmt.Fprint(w, "hello, world\n")
}
func main () {
    http.HandleFunc("/", index)
    http.ListenAndServe(":9091", nil)
}
