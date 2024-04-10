// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"net/http" // Package for building HTTP servers and clients.
)

// routes sets up the application's routes and returns an http.Handler.
func (app *application) routes() http.Handler {
	// Create a new ServeMux, which is an HTTP request multiplexer (router).
	// It matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.
	mux := http.NewServeMux()

	// Create a new file server for serving static files.
	// The file server is wrapped in the limitedFileSystem to prevent directory listing.
	fileServer := http.FileServer(limitedFileSystem{http.Dir(app.config.staticDir)})

	// Register the file server as a handler function for all URLs that start with "/static/".
	// The http.StripPrefix function modifies the request URL's path before the request reaches the file server.
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register handler functions for URL patterns.
	// When a request URL matches one of these patterns, the corresponding handler function is called.
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Wrap the mux (ServeMux) with the recoverPanic, logRequest, and secureHeaders middleware functions.
	// This means that every request will go through these middleware functions in the order they are listed.
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
