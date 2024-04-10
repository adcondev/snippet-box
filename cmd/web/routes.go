package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	// Create a new ServeMux.
	mux := http.NewServeMux()

	// Serve static files from "./ui/static/" directory.
	fileServer := http.FileServer(limitedFileSystem{http.Dir(app.config.staticDir)})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register handler functions for URL patterns.
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return app.logRequest(secureHeaders(mux))
}
