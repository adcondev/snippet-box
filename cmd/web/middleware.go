// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"fmt"      // Package for formatted I/O.
	"net/http" // Package for building HTTP servers and clients.
)

// secureHeaders is a middleware that adds secure headers to the response.
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

// logRequest is a middleware that logs the request.
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request.
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}

// recoverPanic is a middleware that recovers from panic.
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
