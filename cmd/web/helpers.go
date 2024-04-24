// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"bytes" // Package for manipulating byte slices.
	"errors"
	"fmt"      // Package for formatted I/O.
	"net/http" // Package for building HTTP servers and clients.
	// Package for manipulating file paths.
	"runtime/debug" // Package for providing information about the Go runtime.
	"time"          // Package for measuring and displaying time.

	"github.com/go-playground/form/v4"
)

// serverError is a helper function that writes an error message and stack trace to the errorLog,
// then sends a 500 Internal Server Error response to the user. It takes an http.ResponseWriter to
// write the response to, and an error to log and respond with.
func (app *application) serverError(w http.ResponseWriter, err error) {
	// Create a stack trace and store it in the variable trace.
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// Write the error message and stack trace to the errorLog.
	app.errorLog.Output(2, trace)
	// Use the http.Error function to send a 500 status to the user.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError is a helper function that sends a specific status code and corresponding description
// to the user. It takes an http.ResponseWriter to write the response to, and a status code to respond with.
// This function can be used to send any 4xx status code to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	// Use the http.Error function to send the status code and description to the user.
	http.Error(w, http.StatusText(status), status)
}

// notFound is a helper function that sends a 404 Not Found status to the user.
// It uses the clientError function to send the status code and description to the user.
func (app *application) notFound(w http.ResponseWriter) {
	// Use the clientError function to send a 404 status to the user.
	app.clientError(w, http.StatusNotFound)
}

// render is a helper function that renders a template. It writes the rendered template to the
// http.ResponseWriter, along with the provided HTTP status code. If the template does not exist
// in the cache, it sends a server error response. If there's an error when executing the template,
// it also sends a server error response.
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Try to get the template set for the provided page from the cache.
	ts, ok := app.templateCache[page]
	// If the template set is not in the cache, that means the template does not exist.
	// In that case, send a server error response.
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Create a new bytes.Buffer to hold the rendered template.
	// This buffer is an io.Writer, so we can write the rendered template to it.
	buf := new(bytes.Buffer)
	// Render the template and write it to the buffer.
	// If there's an error, send a server error response.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write the HTTP status code to the http.ResponseWriter header.
	w.WriteHeader(status)
	// Write the contents of the buffer to the http.ResponseWriter.
	// This sends the rendered template as the response body.
	buf.WriteTo(w)
}

// newTemplateData is a helper function that creates a new instance of templateData.
// It initializes the CurrentYear field to the current year.
// This function is useful when you need to create a new templateData instance with the CurrentYear field already set.
func (app *application) newTemplateData(r *http.Request) *templateData {
	// Create a new templateData instance.
	// Set the CurrentYear field to the current year.
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func (app *application) decodePostForm(r *http.Request, target any) error {

	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(target, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
