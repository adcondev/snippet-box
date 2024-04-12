// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"net/http" // Package for building HTTP servers and clients.

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// router sets up the application routes. It returns an http.Handler that is the root of the application's routing hierarchy.
// It creates a new httprouter.Router, registers handler functions for URL patterns, and wraps the router with middleware functions.
// This function is useful for centralizing the application's routing logic.
func (app *application) routes() http.Handler {
	// Create a new ServeMux, which is an HTTP request multiplexer (router).
	router := httprouter.New()

	// Register a handler function for the root URL ("/").
	// If the request URL does not match any registered patterns, the NotFoundHandler is called.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Create a new file server for serving static files.
	// The file server is wrapped in the limitedFileSystem to prevent directory listing.
	fileServer := http.FileServer(limitedFileSystem{http.Dir(app.config.StaticDir)})

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

	// Wrap the router with the recoverPanic, logRequest, and secureHeaders middleware functions.
	// This means that every request will go through these middleware functions in the order they are listed.
	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
		secureHeaders,
	)

	// Return the router.
	return standard.Then(router)
}
