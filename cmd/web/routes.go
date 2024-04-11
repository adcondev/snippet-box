// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"net/http" // Package for building HTTP servers and clients.

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// routes sets up the application's routes and returns an http.Handler.
func (app *application) routes() http.Handler {
	// Create a new ServeMux, which is an HTTP request multiplexer (router).
	// It matches the URL of each incoming request against a list of registered patterns and calls the handler for the pattern that most closely matches the URL.
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Create a new file server for serving static files.
	// The file server is wrapped in the limitedFileSystem to prevent directory listing.
	fileServer := http.FileServer(limitedFileSystem{http.Dir(app.config.staticDir)})

	// Register the file server as a handler function for all URLs that start with "/static/".
	// The http.StripPrefix function modifies the request URL's path before the request reaches the file server.
	router.Handler(http.MethodGet, "/static", http.NotFoundHandler())
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Register handler functions for URL patterns.
	// When a request URL matches one of these patterns, the corresponding handler function is called.
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Wrap the mux (ServeMux) with the recoverPanic, logRequest, and secureHeaders middleware functions.
	// This means that every request will go through these middleware functions in the order they are listed.
	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
		secureHeaders,
	)

	return standard.Then(router)
}
