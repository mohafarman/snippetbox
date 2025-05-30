package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	addr := flag.String("port", "4000", "HTTP server port adress.")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
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
	err := server.ListenAndServe()
	errorLog.Print("Server failed to run. Error: " + err.Error())
}
