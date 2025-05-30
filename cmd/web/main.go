package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	mux := http.NewServeMux()

	/* INFO: http.FileServer will transform os.ErrNotExist from
	   neuteredFS.Open() to 404 Not Found response. */
	fs := http.FileServer(neuteredFS{http.Dir("./ui/static/")})
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view/{id}", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	server := http.Server{
		Addr:         ":" + *addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  90 * time.Second,
		ErrorLog:     errorLog,
		Handler:      mux,
	}

	infoLog.Print("Listening on port " + *addr)
	err := server.ListenAndServe()
	errorLog.Print("Server failed to run. Error: " + err.Error())
}

type neuteredFS struct {
	fs http.FileSystem
}

func (nfs neuteredFS) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			/* INFO: Will return an os.ErrNotExist. */
			return nil, err
		}
	}

	return f, nil
}
