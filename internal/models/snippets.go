// Package models contains the application's data models.
package models

// Import the necessary packages.
import (
	"database/sql" // Package for interacting with SQL databases.
	"errors"       // Package for creating error messages.
	"time"         // Package for measuring and displaying time.
)

// Snippet is a type that represents a snippet in the database.
type Snippet struct {
	ID      int       // The ID of the snippet.
	Title   string    // The title of the snippet.
	Content string    // The content of the snippet.
	Created time.Time // The time when the snippet was created.
	Expires time.Time // The time when the snippet expires.
}

// SnippetModel wraps a sql.DB connection pool.
type SnippetModel struct {
	db         *sql.DB   // The database connection pool.
	InsertStmt *sql.Stmt // The prepared statement for inserting a snippet.
	GetStmt    *sql.Stmt // The prepared statement for getting a snippet.
	LatestStmt *sql.Stmt // The prepared statement for getting the latest snippets.
}

// NewSnippetModel creates a new SnippetModel.
func NewSnippetModel(db *sql.DB) (*SnippetModel, error) {

	// Define the SQL for inserting a snippet.
	insert := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Prepare the SQL statement.
	// If there's an error (for example, if the SQL statement is invalid), return nil and the error.
	insertStmt, err := db.Prepare(insert)
	if err != nil {
		return nil, err
	}

	// Define the SQL for getting a snippet.
	get := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// Prepare the SQL statement.
	// If there's an error (for example, if the SQL statement is invalid), return nil and the error.
	getStmt, err := db.Prepare(get)
	if err != nil {
		return nil, err
	}

	// Define the SQL for getting the latest snippets.
	latest := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// Prepare the SQL statement.
	// If there's an error (for example, if the SQL statement is invalid), return nil and the error.
	latestStmt, err := db.Prepare(latest)
	if err != nil {
		return nil, err
	}

	// Return a new SnippetModel with the database connection and the prepared statements.
	return &SnippetModel{db, insertStmt, getStmt, latestStmt}, nil
}

// Insert inserts a new snippet into the database.
func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	// Start a new transaction.
	// If there's an error (for example, if the transaction can't be started), return 0 and the error.
	tx, err := sm.db.Begin()
	if err != nil {
		return 0, err
	}

	// Use the defer keyword to ensure that the transaction is rolled back if any subsequent code returns an error.
	defer tx.Rollback()

	// Execute the prepared statement for inserting a snippet.
	// If there's an error (for example, if the SQL statement is invalid), return 0 and the error.
	res, err := tx.Stmt(sm.InsertStmt).Exec(title, content, expires)
	if err != nil {
		return 0, err
	}

	// Commit the transaction.
	// If there's an error (for example, if the transaction can't be committed), return 0 and the error.
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	// Get the ID of the new snippet.
	// If there's an error (for example, if the ID can't be retrieved), return 0 and the error.
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// If there's no error, return the ID of the new snippet and nil for the error.
	return int(id), nil
}

// Get retrieves a snippet from the database based on its ID.
func (sm *SnippetModel) Get(id int) (*Snippet, error) {

	// Create a new Snippet struct.
	s := &Snippet{}

	// Execute the prepared statement for getting a snippet.
	// Scan the result into the Snippet struct.
	// If there's an error (for example, if the SQL statement is invalid), handle it in the next block.
	err := sm.GetStmt.QueryRow(id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	// If there's an error...
	if err != nil {
		// If the error is that no rows were returned from the query, return nil and the ErrNoRecord error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			// If it's a different error, return nil and the error.
			return nil, err
		}
	}

	// If there's no error, return the Snippet struct and nil for the error.
	return s, nil
}

// Latest retrieves the latest snippets from the database.
func (sm *SnippetModel) Latest() ([]*Snippet, error) {

	// Execute the prepared statement for getting the latest snippets.
	// If there's an error (for example, if the SQL statement is invalid), return nil and the error.
	rows, err := sm.LatestStmt.Query()
	if err != nil {
		return nil, err
	}
	// Use the defer keyword to ensure that the rows are closed at the end, even if an error occurs.
	defer rows.Close()

	// Create a new slice to hold the Snippet structs.
	snippets := []*Snippet{}

	// Loop over the rows.
	for rows.Next() {
		// For each row, create a new Snippet struct.
		s := &Snippet{}
		// Scan the row into the Snippet struct.
		// If there's an error (for example, if the row can't be scanned), return nil and the error.
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append the Snippet struct to the slice.
		snippets = append(snippets, s)
	}
	// If there's an error with the rows (for example, if there's a problem with the iteration), return nil and the error.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If there's no error, return the slice of Snippet structs and nil for the error.
	return snippets, nil
}
