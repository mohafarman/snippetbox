package main

import (
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mohafarman/snippetbox/internal/models"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	snippets *models.SnippetModel
}

func main() {
	addr := flag.String("port", "4000", "HTTP server port adress.")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	/* INFO: driver-specific parameter which instructs
	our driver to convert SQL TIME and DATE fields to Go time.Time objects */
	dsn := "snippetbox?parseTime=true"
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{
			DB: db,
		},
	}

	server := http.Server{
		Addr:         ":" + *addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  90 * time.Second,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
	}

	infoLog.Print("Listening on port " + *addr)
	err = server.ListenAndServe()
	errorLog.Print("Server failed to run. Error: " + err.Error())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
