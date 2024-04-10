// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"bytes"         // Package for manipulating byte slices.
	"fmt"           // Package for formatted I/O.
	"net/http"      // Package for building HTTP servers and clients.
	"path/filepath" // Package for manipulating file paths.
	"runtime/debug" // Package for providing information about the Go runtime.
	"time"          // Package for measuring and displaying time.
)

// limitedFileSystem wraps an http.FileSystem to disable directory listings.
type limitedFileSystem struct {
	fs http.FileSystem // The underlying file system.
}

// Open opens a file in the limitedFileSystem.
func (nfs limitedFileSystem) Open(path string) (http.File, error) {
	// Open the file.
	f, err := nfs.fs.Open(path)
	// If there's an error (for example, the file doesn't exist), return it.
	if err != nil {
		return nil, err
	}

	// Get the file's metadata.
	s, _ := f.Stat()
	// If the file is a directory...
	if s.IsDir() {
		// ...try to open its index.html file.
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			// If there's an error (which means the index.html file doesn't exist)...
			// ...close the directory...
			closeErr := f.Close()
			// ...and if there's an error when closing the directory, return it.
			if closeErr != nil {
				return nil, closeErr
			}
			// Otherwise, return the original error.
			return nil, err
		}
	}

	// If there's no error, return the file.
	return f, nil
}

// serverError is a helper function that writes an error message and stack trace to the errorLog,
// then sends a 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	// Create a stack trace and store it in the variable trace.
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// Write the error message and stack trace to the errorLog.
	app.errorLog.Output(2, trace)
	// Use the http.Error function to send a 500 status to the user.
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError is a helper function that sends a specific status code and corresponding description
// to the user. You can use this function to send any 4xx status code to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	// Use the http.Error function to send the status code and description to the user.
	http.Error(w, http.StatusText(status), status)
}

// notFound is a helper function that sends a 404 Not Found status to the user.
func (app *application) notFound(w http.ResponseWriter) {
	// Use the clientError function to send a 404 status to the user.
	app.clientError(w, http.StatusNotFound)
}

// render is a helper function that renders a template. It writes the rendered template to the
// http.ResponseWriter, along with the provided HTTP status code.
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

// newTemplateData is a helper function that creates a new templateData struct and initializes
// its CurrentYear field to the current year.
func (app *application) newTemplateData() *templateData {
	// Create a new templateData struct.
	// Initialize the CurrentYear field to the current year.
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
