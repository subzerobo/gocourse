package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	// http.DefaultServeMux
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "This is my simple web-server")
	})

	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		fragment := r.URL.Path
		fmt.Println(fragment)
		//values := r.URL.Query()
		//values.Get("name")

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 3 {
			http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
			return
		}
		name := parts[2]

		fmt.Fprintf(w, "Hello %s\n", name)
	})

	http.ListenAndServe(":8080", nil)

}

func main2() {
	// http.DefaultServeMux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Got / request\n")
		io.WriteString(w, "This is my website!\n")
	})

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("got /hello request\n")
		io.WriteString(w, "Hello, HTTP!\n")
	})

	httpServer := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	httpServer.ListenAndServe()
}
