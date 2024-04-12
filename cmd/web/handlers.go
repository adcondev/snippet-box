// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"errors"   // Package for creating error messages.
	"fmt"      // Package for formatted I/O.
	"net/http" // Package for building HTTP servers and clients.
	"strconv"  // Package for converting strings to numeric types.
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox.consdotpy.xyz/internal/models" // Import the models package.
)

type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

// home is a handler function that serves the root URL ("/").
func (app *application) home(w http.ResponseWriter, r *http.Request) {
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

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.html", data)
}

// snippetCreate is a handler function that serves the "/snippet/create" URL.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	if expires != 1 && expires != 7 && expires != 365 {
		form.FieldErrors["expires"] = "this fields must equal 1, 7  or 365"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	// Insert the new snippet into the database.
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	// If there's an error (for example, a database error), send a server error response.
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If there's no error, the snippet was inserted successfully.
	// Redirect the client to the page for the new snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
