// Package models contains the application's data models.
package models

// Import the necessary packages.
import (
	"database/sql" // Package for interacting with SQL databases.
	"errors"       // Package for creating error messages.
	"time"         // Package for measuring and displaying time.
)

// Snippet represents a snippet in the application. It is used to hold data related to a snippet.
// A snippet consists of an ID, a title, content, and timestamps for when the snippet was created and when it expires.
type Snippet struct {
	ID      int       // ID is the unique identifier for the snippet.
	Title   string    // Title is the title of the snippet.
	Content string    // Content is the content of the snippet.
	Created time.Time // Created is the time when the snippet was created.
	Expires time.Time // Expires is the time when the snippet expires.
}

// SnippetModel wraps a sql.DB connection pool and provides methods for interacting with the snippets table in the database.
// It holds prepared SQL statements for inserting a snippet, getting a snippet, and getting the latest snippets.
// This struct is useful for encapsulating the database operations related to snippets.
type SnippetModel struct {
	DB         *sql.DB   // DB is the database connection pool.
	InsertStmt *sql.Stmt // InsertStmt is the prepared statement for inserting a snippet.
	GetStmt    *sql.Stmt // GetStmt is the prepared statement for getting a snippet.
	LatestStmt *sql.Stmt // LatestStmt is the prepared statement for getting the latest snippets.
}

// NewSnippetModel creates a new SnippetModel with a given database connection.
// It prepares SQL statements for inserting a snippet, getting a snippet, and getting the latest snippets.
// These prepared statements are stored in the SnippetModel, which can then be used to perform these operations.
// This function is useful for setting up the SnippetModel with the SQL statements it needs to interact with the database.
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

// Insert inserts a new snippet into the database. It starts a new transaction, executes the prepared statement for inserting a snippet,
// commits the transaction, and retrieves the ID of the new snippet. If there's an error at any point (for example, if the transaction can't be started,
// if the SQL statement is invalid, if the transaction can't be committed, or if the ID can't be retrieved), it returns 0 and the error.
// If there's no error, it returns the ID of the new snippet and nil for the error.
func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	// Start a new transaction.
	// If there's an error (for example, if the transaction can't be started), return 0 and the error.
	tx, err := sm.DB.Begin()
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

// Get retrieves a snippet from the database based on its ID. It executes the prepared statement for getting a snippet,
// and scans the result into a new Snippet struct. If there's an error (for example, if the SQL statement is invalid),
// it handles it accordingly: if the error is that no rows were returned from the query, it returns nil and the ErrNoRecord error;
// if it's a different error, it returns nil and the error. If there's no error, it returns the Snippet struct and nil for the error.
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

// Latest retrieves the 10 most recently created snippets that have not expired from the database. It executes the prepared statement for getting the latest snippets,
// and scans the results into a slice of Snippet structs. If there's an error (for example, if the SQL statement is invalid),
// it returns nil and the error. If there's no error, it returns the slice of Snippet structs and nil for the error.
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
