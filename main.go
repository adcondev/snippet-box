package main

import (
	"log"
	"net/http"
)

// home is a handler function that serves the root URL ("/").
func home(w http.ResponseWriter, r *http.Request) {
	// If the requested URL path isn't exactly "/", then respond with a 404 status.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Otherwise, respond with a "Hello from Snippetbox" message.
	w.Write([]byte("Hello from Snippetbox"))
}

// snippetView is a handler function that serves the "/snippet/view" URL.
func snippetView(w http.ResponseWriter, r *http.Request) {
	// For now, it simply responds with a static message.
	w.Write([]byte("Display a specific snippet..."))
}

// snippetCreate is a handler function that serves the "/snippet/create" URL.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// For now, it simply responds with a static message.
	w.Write([]byte("Create a new snippet..."))
}

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
