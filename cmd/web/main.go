package main

import (
	"flag"
	"log/slog"
	"net/http"
	"path/filepath"
	"time"
)

func main() {

	addr := flag.String("port", "4000", "HTTP server port adress.")
	flag.Parse()

	mux := http.NewServeMux()

	/* INFO: http.FileServer will transform os.ErrNotExist from
	   neuteredFS.Open() to 404 Not Found response. */
	fs := http.FileServer(neuteredFS{http.Dir("./ui/static/")})
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view/{id}", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	server := http.Server{
		Addr:         ":" + *addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  90 * time.Second,
		Handler:      mux,
	}

	slog.Info("Listening on port " + *addr)
	err := server.ListenAndServe()
	slog.Error("Server failed to run. Error: " + err.Error())
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
