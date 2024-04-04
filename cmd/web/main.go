package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type configuration struct {
	addr      string
	staticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config   configuration
}

// main is the application's entry point.
func main() {
	var cfg configuration
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static/", "Path to static assets")
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

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		config:   cfg,
	}

	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	// Start server on port 4000.
	infoLog.Printf("Starting server on %s", cfg.addr)
	err := srv.ListenAndServe()

	// Log and exit on server start error.
	errorLog.Fatal(err)
}
