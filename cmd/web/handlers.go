package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

// home is a handler function that serves the root URL ("/").
func home(w http.ResponseWriter, r *http.Request) {
	// Clean the requested URL path to prevent directory traversal attacks.
	r.URL.Path = path.Clean(r.URL.Path)

	// If the cleaned URL path isn't exactly "/", then respond with a 404 status.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Define a slice that holds the file paths for the templates.
	files := []string{
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
		"./ui/html/base.html",
	}

	// Use the template.ParseFiles function to read the template files and store the templates in a template set.
	ts, err := template.ParseFiles(files...)
	// If there's an error, log the detailed error message and send a generic 500 Internal Server Error response to the user.
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Use the ExecuteTemplate method to write the "base" template to the http.ResponseWriter.
	// For now we're passing in nil as the last parameter, because we're not displaying any dynamic data.
	err = ts.ExecuteTemplate(w, "base", nil)
	// If there's an error, log the detailed error message and send a generic 500 Internal Server Error response to the user.
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// snippetView is a handler function that serves the "/snippet/view" URL.
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Convert the id from the URL query to an integer.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// If the conversion fails or the id is less than 1, respond with a 404 status.
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// For now, it simply responds with a static message.
	fmt.Fprintf(w, "Display a specific snippet with ID %d", id)
}

// snippetCreate is a handler function that serves the "/snippet/create" URL.
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// If the request method is not POST, respond with a "Method Not Allowed" status.
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Set the response content type to JSON.
	w.Header().Set("Content-Type", "application/json")
	// For now, it simply responds with a static JSON message.
	w.Write([]byte(`{"Create":"Snippet"}`))
}

// downloadSnippet is a handler function that serves the "/download/" URL.
func downloadSnippet(w http.ResponseWriter, r *http.Request) {
	// Clean the URL path to prevent directory traversal attacks.
	r.URL.Path = path.Clean(r.URL.Path)
	// Serve the static file.
	http.ServeFile(w, r, "./ui/static/img/logo.png")
}
