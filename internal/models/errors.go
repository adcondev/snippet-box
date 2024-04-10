// Package models contains the application's data models.
package models

// Import the errors package to create error messages.
import (
	"errors"
)

// ErrNoRecord is an error that is returned when a database query returns no results.
// It's created using the errors.New function, which creates a new error with the specified message.
var ErrNoRecord = errors.New("models: no matching record found")
