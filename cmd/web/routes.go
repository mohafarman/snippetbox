package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	/* INFO: http.FileServer will transform os.ErrNotExist from
	   neuteredFS.Open() to 404 Not Found response. */
	fs := http.FileServer(neuteredFS{http.Dir("./ui/static/")})
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view/{id}", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
