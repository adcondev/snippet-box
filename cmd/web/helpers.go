package main

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"
	"time"
)

// limitedFileSystem wraps an http.FileSystem to disable directory listings.
type limitedFileSystem struct {
	fs http.FileSystem
}

// Open opens a file in the limitedFileSystem.
// If the file is a directory, Open tries to open its index.html file.
// If the index.html file doesn't exist, Open returns an error.
func (nfs limitedFileSystem) Open(path string) (http.File, error) {
	// Open the file.
	f, err := nfs.fs.Open(path)
	if err != nil {
		// If there's an error, return it.
		return nil, err
	}

	// Get the file's metadata.
	s, _ := f.Stat()
	if s.IsDir() {
		// If the file is a directory, try to open its index.html file.
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			// If there's an error (which means the index.html file doesn't exist), close the directory and return the error.
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	// Return the file.
	return f, nil
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(
		w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError,
	)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}

func (app *application) newTemplateData() *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
