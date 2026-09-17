package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"Golang/api"
	"Golang/api/yqk"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/yqk/", func(w http.ResponseWriter, r *http.Request) {
		yqk.Handler(w, r)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.HasPrefix(path, "/live/") {
			api.LiveHandler(w, r)
			return
		}
		api.Handler(w, r)
	})

	addr := ":" + port
	fmt.Printf("Server listening on http://0.0.0.0%s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}