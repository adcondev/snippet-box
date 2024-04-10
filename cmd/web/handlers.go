// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"errors"   // Package for creating error messages.
	"fmt"      // Package for formatted I/O.
	"net/http" // Package for building HTTP servers and clients.
	"strconv"  // Package for converting strings to numeric types.

	"snippetbox.consdotpy.xyz/internal/models" // Import the models package.
)

// home is a handler function that serves the root URL ("/").
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// If the URL path isn't exactly "/", respond with a 404 status.
	// This means the function only responds to the exact path "/" and not any other path.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Fetch the latest snippets from the database.
	// The Latest method is expected to return the most recent snippets.
	snippets, err := app.snippets.Latest()
	// If there's an error (for example, a database error), send a server error response.
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Create a new template data map and add the snippets to it.
	// This map will be passed to the template for rendering.
	data := app.newTemplateData()
	data.SnippetsData = snippets

	// Render the home page with the snippets.
	// The render method is expected to render the "home.html" template with the provided data.
	app.render(w, http.StatusOK, "home.html", data)
}

// snippetView is a handler function that serves the "/snippet/view" URL.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Convert the id from the URL query to an integer.
	// The id is expected to be passed in the query string like "/snippet/view?id=1".
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// If the conversion fails (which means the id is not a valid integer) or the id is less than 1,
	// respond with a 404 status by calling the notFound helper.
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Fetch the snippet with the given id from the database.
	snippet, err := app.snippets.Get(id)
	// If there's an error, handle it.
	if err != nil {
		// If the error is of type models.ErrNoRecord, that means no snippet with the given id was found,
		// so respond with a 404 status by calling the notFound helper.
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			// For any other kind of error, respond with a 500 status by calling the serverError helper.
			app.serverError(w, err)
		}
		return
	}

	// If there's no error, the snippet was fetched successfully.
	// Create a new template data map and add the snippet to it.
	// This map will be passed to the template for rendering.
	data := app.newTemplateData()
	data.SnippetData = snippet

	// Render the "view.html" template with the provided data.
	app.render(w, http.StatusOK, "view.html", data)
}

// snippetCreate is a handler function that serves the "/snippet/create" URL.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// If the request method is not POST, respond with a "Method Not Allowed" status.
	// This means the function only responds to POST requests and not any other method.
	if r.Method != http.MethodPost {
		// Set the "Allow" header to "POST" to tell the client that only POST requests are allowed.
		w.Header().Set("Allow", http.MethodPost)
		// Respond with a 405 status by calling the clientError helper.
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Define the title, content, and expiration for the new snippet.
	// In a real application, these would probably come from the request body.
	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// Insert the new snippet into the database.
	id, err := app.snippets.Insert(title, content, expires)
	// If there's an error (for example, a database error), send a server error response.
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If there's no error, the snippet was inserted successfully.
	// Redirect the client to the page for the new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
