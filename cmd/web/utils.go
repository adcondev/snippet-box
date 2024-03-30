// Package main provides the main entry point to the application.
package main

import (
	"net/http"
	"path/filepath"
)

// neuteredFileSystem wraps an http.FileSystem to disable directory listings.
type neuteredFileSystem struct {
	fs http.FileSystem
}

// Open opens a file in the neuteredFileSystem.
// If the file is a directory, Open tries to open its index.html file.
// If the index.html file doesn't exist, Open returns an error.
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
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
