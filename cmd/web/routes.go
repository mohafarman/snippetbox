package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	/* INFO: http.FileServer will transform os.ErrNotExist from
	   neuteredFS.Open() to 404 Not Found response. */
	fs := http.FileServer(neuteredFS{http.Dir("./ui/static/")})
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	/* Using "GET /" causes panic:
	pattern "GET /" (...) conflicts with pattern "/static/" (...) */
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	/* INFO: flow of exeuction:
	   secureHeaders → servemux → application handler → servemux → secureHeaders */
	// INFO: Without alice: return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	return standard.Then(mux)
}
