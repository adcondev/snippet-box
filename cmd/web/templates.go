// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"path/filepath" // Package for manipulating file paths.
	"text/template" // Package for manipulating text templates.
	"time"          // Package for measuring and displaying time.

	"snippetbox.consdotpy.xyz/internal/models" // Import the models package.
)

// templateData holds data to be passed into templates.
type templateData struct {
	CurrentYear  int               // The current year.
	SnippetData  *models.Snippet   // Data for a single snippet.
	SnippetsData []*models.Snippet // Data for multiple snippets.
	Form         any
}

// functions is a map that acts as a lookup for functions that can be used in templates.
var functions = template.FuncMap{
	"humanDate": humanDate, // Map the "humanDate" key to the humanDate function.
}

// humanDate formats a time.Time object to a human-friendly date format.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// newTemplateCache creates a new template cache as a map and returns it.
func newTemplateCache() (map[string]*template.Template, error) {

	// Create a new template cache.
	cache := map[string]*template.Template{}

	// Get a slice of all filepaths with the .html extension in the ui/html/pages folder.
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	// If there's an error, return the cache and the error.
	if err != nil {
		return nil, err
	}

	// Loop over the pages.
	for _, page := range pages {
		// Extract the file name (like 'home.page.html') from the full file path and assign it to the name variable.
		name := filepath.Base(page)

		// Create a new template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// Add the partials layouts to the template set.
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// Parse templates in pages folder
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Store the template set in the cache, using the page name (like 'home.page.html') as the key.
		cache[name] = ts
	}

	// Return the template cache.
	return cache, nil
}
