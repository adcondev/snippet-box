// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"context"
	"fmt"      // Package for formatted I/O.
	"net/http" // Package for building HTTP servers and clients.
)

// secureHeaders is a middleware function that adds secure headers to the HTTP response.
// It takes an http.Handler as input and returns an http.Handler.
// The returned http.Handler adds several secure headers to the response header and then calls the ServeHTTP method of the input handler.
// This function is useful for adding secure headers to all responses in a centralized way.
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add secure headers to the response.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// logRequest is a middleware function that logs the details of each HTTP request.
// It takes an http.Handler as input and returns an http.Handler.
// The returned http.Handler logs the remote address, protocol, method, and URL of the request, and then calls the ServeHTTP method of the input handler.
// This function is useful for logging the details of each request in a centralized way.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the remote address, protocol, method, and URL of the request.
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// recoverPanic is a middleware function that recovers from any panics and writes a 500 Internal Server Error response.
// It takes an http.Handler as input and returns an http.Handler.
// The returned http.Handler uses the defer keyword to ensure that the function is called at the end, even if a panic occurs.
// If a panic occurs, it sets the connection header to "close", logs the error, and sends a 500 Internal Server Error response.
// This function is useful for recovering from panics in a centralized way and providing a user-friendly error message.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use the defer keyword to ensure that this function is called at the end, even if a panic occurs.
		defer func() {
			// Use the recover function to catch a panic.
			if err := recover(); err != nil {
				// If a panic occurred, set the connection header to "close".
				w.Header().Set("Connection", "close")
				// Log the error and send a 500 Internal Server Error response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.users.Exists(id)
		if err != nil {
			app.serverError(w, err)
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
