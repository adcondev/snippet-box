package main

import (
	"log"
	"net/http"
)

// main is the entry point of the application.
func main() {
	// Create a new ServeMux.
	mux := http.NewServeMux()
	// Register the handler functions with the ServeMux for their respective URL patterns.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Log a message to indicate that the server is starting.
	log.Print("Starting server on :4000")
	// Start the web server.
	err := http.ListenAndServe(":4000", mux)
	// If http.ListenAndServe returns an error, log the error and exit the program.
	log.Fatal(err)
}
