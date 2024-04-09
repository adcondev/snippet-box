package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.consdotpy.xyz/internal/models"
)

// home is a handler function that serves the root URL ("/").
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// If the cleaned URL path isn't exactly "/", then respond with a 404 status.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.SnippetsData = snippets

	app.render(w, http.StatusOK, "home.html", data)
}

// snippetView is a handler function that serves the "/snippet/view" URL.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Convert the id from the URL query to an integer.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// If the conversion fails or the id is less than 1, respond with a 404 status.
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData()
	data.SnippetData = snippet

	app.render(w, http.StatusOK, "view.html", data)

}

// snippetCreate is a handler function that serves the "/snippet/create" URL.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// If the request method is not POST, respond with a "Method Not Allowed" status.
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
