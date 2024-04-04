package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

	infoLog := log.New(
		os.Stdout,
		"INFO\t",
		log.Ldate|log.Ltime|log.LUTC,
	)
	errorLog := log.New(
		os.Stderr,
		"ERROR\t",
		log.Ldate|log.Ltime|log.LUTC|log.Llongfile,
	)

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

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Start server on port 4000.
	infoLog.Printf("Starting server on %s", cfg.addr)
	err := srv.ListenAndServe()

	// Log and exit on server start error.
	errorLog.Fatal(err)
}
