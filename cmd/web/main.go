package main

import (
	"flag"
	"log"
	"net/http"
)

type config struct {
	addr      string
	staticDir string
}

// main is the application's entry point.
func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.Parse()
	// Create a new ServeMux.
	mux := http.NewServeMux()

	// Serve static files from "./ui/static/" directory.
	fileServer := http.FileServer(limitedFileSystem{http.Dir(cfg.staticDir)})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register handler functions for URL patterns.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	mux.HandleFunc("/download/", snippetDownload)

	// Start server on port 4000.
	log.Printf("Starting server on %s", cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)

	// Log and exit on server start error.
	log.Fatal(err)
}
