// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"net/http" // Package for building HTTP servers and clients.

	"snippetbox.adcon.dev/ui"

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

	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	router.HandlerFunc(http.MethodGet, "/ping", ping)

	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticate)

	// Register handler functions for URL patterns.
	// When a request URL matches one of these patterns, the corresponding handler function is called.
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

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
