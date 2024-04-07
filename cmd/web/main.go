package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"snippetbox.consdotpy.xyz/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type configuration struct {
	addr      string
	staticDir string
	dsn       string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config   configuration
	snippets *models.SnippetModel
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// main is the application's entry point.
func main() {
	var cfg configuration
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(
		os.Stdout,
		"INFO\t",
		log.Ldate|log.Ltime|log.LUTC,
	)
	errorLog := log.New(
		os.Stderr,
		"ERROR\t",
		log.Ldate|log.Ltime|log.LUTC|log.Llongfile,
	)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		config:   cfg,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Start server.
	infoLog.Printf("Starting server on %s", cfg.addr)
	err = srv.ListenAndServe()

	// Log and exit on server start error.
	errorLog.Fatal(err)
}
