package main

import (
	"log"
	"net/http"
	"path"
)

// cleanFileServer is a wrapper around http.FileServer that cleans the URL path before serving files.
func cleanFileServer(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = path.Clean(r.URL.Path)
		fs.ServeHTTP(w, r)
	}
}

// main is the application's entry point.
func main() {
	// Create a new ServeMux.
	mux := http.NewServeMux()

	// Serve static files from "./ui/static/" directory.
	fileServer := cleanFileServer(http.FileServer(http.Dir("./ui/static/")))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register handler functions for URL patterns.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	mux.HandleFunc("/download/", downloadSnippet)

	// Start server on port 4000.
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)

	// Log and exit on server start error.
	log.Fatal(err)
}
