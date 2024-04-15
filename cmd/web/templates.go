// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"path/filepath" // Package for manipulating file paths.
	"text/template" // Package for manipulating text templates.
	"time"          // Package for measuring and displaying time.

	"snippetbox.consdotpy.xyz/internal/models" // Import the models package.
)

// templateData holds data to be passed into templates. It is used to provide a consistent
// structure for passing data to templates, making it easier to manage and evolve over time.
type templateData struct {
	CurrentYear  int               // CurrentYear holds the current year.
	SnippetData  *models.Snippet   // SnippetData holds data for a single snippet.
	SnippetsData []*models.Snippet // SnippetsData holds data for multiple snippets.
	Form         any               // Form holds form data.
	Flash        string
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
// The cache is a map where the keys are page names (like 'home.page.html') and the values are the corresponding templates.
// This function is useful for preloading all the templates into the cache on application startup.
// This means that the templates do not need to be loaded from the disk every time a request is made, which improves the performance of the application.
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
