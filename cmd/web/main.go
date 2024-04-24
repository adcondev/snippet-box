// Package main is the main package for this application.
package main

// Import the necessary packages.
import (
	"crypto/tls"
	"database/sql"  // Package for interacting with SQL databases.
	"flag"          // Package for parsing command-line flags.
	"log"           // Package for logging.
	"net/http"      // Package for building HTTP servers and clients.
	"os"            // Package for interacting with the operating system.
	"text/template" // Package for manipulating text templates.
	"time"

	"snippetbox.adcon.dev/internal/models" // Import the models package.

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver.
)

// configuration represents the application configuration. It includes fields for each configuration option.
// These fields are populated with values from environment variables when the application starts.
// This struct is useful for centralizing all configuration options and making them available throughout the application.
type configuration struct {
	Addr      string // Addr is the network address that the application should listen on.
	StaticDir string // StaticDir is the directory where static files are stored.
	Dsn       string // Secret is the secret key used for session authentication.
}

// application holds the application-wide dependencies. It includes fields for the error and info loggers,
// the application configuration, the model for snippets, and the cache for templates.
// This struct is useful for making these dependencies available throughout the application.
type application struct {
	errorLog       *log.Logger                   // errorLog is the logger for errors.
	infoLog        *log.Logger                   // infoLog is the logger for information.
	config         configuration                 // config is the application configuration.
	snippets       *models.SnippetModel          // snippets is the model for snippets.
	templateCache  map[string]*template.Template // templateCache is the cache for templates.
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	users          *models.UserModel
}

// openDB opens a new database connection with the provided data source name (DSN).
// It uses the sql.Open function to open a new database connection and the db.Ping function to establish a connection
// and verify that the given DSN is valid. If there's an error when opening the connection or when pinging the database,
// it returns nil and the error. If there's no error, it returns the database connection and nil for the error.
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

// main is the application's entry point. It sets up the application configuration, loggers, database connection,
// and HTTP server. It also handles any errors that occur during setup.
func main() {
	// Create a new configuration struct and parse command-line flags into it.
	// The configuration includes the network address, static assets directory, and MySQL data source name.
	var config configuration
	flag.StringVar(&config.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&config.StaticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.StringVar(&config.Dsn, "dsn", "", "MySQL data source name")
	flag.Parse()

	// Create a new logger for informational messages and write them to os.Stdout.
	infoLog := log.New(
		os.Stdout,
		"INFO\t",
		log.Ldate|log.Ltime|log.LUTC,
	)

	// Create a new logger for error messages, write them to os.Stderr, and include more detailed information.
	errorLog := log.New(
		os.Stderr,
		"ERROR\t",
		log.Ldate|log.Ltime|log.LUTC|log.Llongfile,
	)

	// Call the openDB function to open a new database connection.
	db, err := openDB(config.Dsn)
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

	users, err := models.NewUserModel(db)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer users.InsertStmt.Close()
	defer users.AuthStmt.Close()
	defer users.ExistsStmt.Close()

	formDecoder := form.NewDecoder()

	// Call the newTemplateCache function to create a new template cache.
	templateCache, err := newTemplateCache()
	// If there's an error, log the error message and stop the application.
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Create a new application struct and assign the loggers, configuration, snippets model, and template cache.
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		config:         config,
		snippets:       snippets,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		users:          users,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Create a new HTTP server with the network address from the configuration, the error logger, and the application's routes as the handler.
	srv := &http.Server{
		Addr:           config.Addr,
		ErrorLog:       errorLog,
		Handler:        app.routes(),
		TLSConfig:      tlsConfig,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288,
	}

	// Log a message to indicate that the server is starting.
	infoLog.Printf("Starting server on %s", config.Addr)
	// Start the server and listen for requests.
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

	// If there's an error (for example, if the server can't start), log the error message and stop the application.
	errorLog.Fatal(err)
}
