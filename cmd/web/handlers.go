package main

import (
	"fmt"
	"net/http"
	"strconv"
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// For now, it simply responds with a static message.
	fmt.Fprintf(w, "Display a specific snippet with ID %d", id)
}

// snippetCreate is a handler function that serves the "/snippet/create" URL.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	// For now, it simply responds with a static message.
	w.Write([]byte(`{"Create":"Snippet"}`))
}
