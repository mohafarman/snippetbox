package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"
)

func (app *application) errorServer(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) errorClient(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) errorNotFound(w http.ResponseWriter) {
	app.errorClient(w, http.StatusNotFound)
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
