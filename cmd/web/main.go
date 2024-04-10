// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"database/sql"  // Package for interacting with SQL databases.
	"flag"          // Package for parsing command-line flags.
	"log"           // Package for logging.
	"net/http"      // Package for building HTTP servers and clients.
	"os"            // Package for interacting with the operating system.
	"text/template" // Package for manipulating text templates.

	"snippetbox.consdotpy.xyz/internal/models" // Import the models package.

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver.
)

// configuration holds the configuration for the application.
type configuration struct {
	addr      string // The network address to listen on.
	staticDir string // The directory where static files are stored.
	dsn       string // The data source name (DSN) for the database.
}

// application holds the application-wide dependencies.
type application struct {
	errorLog      *log.Logger                   // The logger for errors.
	infoLog       *log.Logger                   // The logger for information.
	config        configuration                 // The application configuration.
	snippets      *models.SnippetModel          // The model for snippets.
	templateCache map[string]*template.Template // The cache for templates.
}

// openDB opens a new database connection with the provided data source name (DSN).
func openDB(dsn string) (*sql.DB, error) {
	// Open a new database connection with the provided DSN.
	// sql.Open does not establish any connections to the database, nor does it validate driver connection parameters.
	db, err := sql.Open("mysql", dsn)
	// If there's an error, return nil and the error.
	if err != nil {
		return nil, err
	}

	// Ping the database to establish a connection and verify that the given DSN is valid.
	if err = db.Ping(); err != nil {
		// If there's an error, return nil and the error.
		return nil, err
	}

	// If there's no error, return the database connection and nil for the error.
	return db, nil
}

// main is the application's entry point.
func main() {
	// Create a new configuration struct.
	var config configuration

	// Use the flag package to define command-line flags for the network address, static assets directory, and MySQL data source name.
	// The flag package will parse the command-line arguments and assign the values to the fields in the config struct.
	flag.StringVar(&config.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&config.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.StringVar(&config.dsn, "dsn", "", "MySQL data source name")
	flag.Parse() // Parse the command-line flags.

	// Create a new logger for informational messages and write them to os.Stdout.
	infoLog := log.New(
		os.Stdout,
		"INFO\t",
		log.Ldate|log.Ltime|log.LUTC,
	)

	// Create a new logger for error messages and write them to os.Stderr.
	// Include more detailed information in error logs.
	errorLog := log.New(
		os.Stderr,
		"ERROR\t",
		log.Ldate|log.Ltime|log.LUTC|log.Llongfile,
	)

	// Call the openDB function to open a new database connection.
	db, err := openDB(config.dsn)
	// If there's an error, log the error message and stop the application.
	if err != nil {
		errorLog.Fatal(err)
	}

	// Close the database connection when the main function exits.
	defer db.Close()

	// Call the NewSnippetModel function to create a new SnippetModel.
	snippets, err := models.NewSnippetModel(db)
	// If there's an error (for example, if the SnippetModel can't be created), log the error message and stop the application.
	if err != nil {
		errorLog.Fatal(err)
	}

	// Close the prepared statements when the main function exits.
	defer snippets.InsertStmt.Close()
	defer snippets.GetStmt.Close()
	defer snippets.LatestStmt.Close()

	// Call the newTemplateCache function to create a new template cache.
	templateCache, err := newTemplateCache()
	// If there's an error, log the error message and stop the application.
	if err != nil {
		errorLog.Fatal(err)
	}

	// Create a new application struct and assign the loggers, configuration, snippets model, and template cache.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		config:        config,
		snippets:      snippets,
		templateCache: templateCache,
	}

	// Create a new HTTP server with the network address from the configuration, the error logger, and the application's routes as the handler.
	srv := &http.Server{
		Addr:     config.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Log a message to indicate that the server is starting.
	infoLog.Printf("Starting server on %s", config.addr)
	// Start the server and listen for requests.
	err = srv.ListenAndServe()

	// If there's an error (for example, if the server can't start), log the error message and stop the application.
	errorLog.Fatal(err)
}
