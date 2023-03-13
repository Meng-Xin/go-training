package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	// http.Handle("/foo", fooHandler)

	http.HandleFunc("/bar", barRes)
	http.HandleFunc("/ping", pingRes)
	http.HandleFunc("/watch", watch)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pingRes(w http.ResponseWriter, r *http.Request) {
	info := "Pong"
	w.Write([]byte(info))
}

func barRes(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func watch(w http.ResponseWriter, r *http.Request) {
	info := "This is Watch"
	w.Write([]byte(info))
}
