package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/mohafarman/snippetbox/ui"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.errorNotFound(w)
	})

	/* INFO: http.FileServer will transform os.ErrNotExist from
	   neuteredFS.Open() to 404 Not Found response. */
	// fs := http.FileServer(neuteredFS{http.Dir("./ui/static/")})
	// router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fs))

	router.HandlerFunc(http.MethodGet, "/ping", app.ping)

	fs := http.FileServer(http.FS(ui.Files))
	// INFO: No need for strip prefix when using embedded fs
	router.Handler(http.MethodGet, "/static/*filepath", fs)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authentication)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	/* INFO: flow of exeuction:
	   secureHeaders → servemux → application handler → servemux → secureHeaders */
	// INFO: Without alice: return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	return standard.Then(router)
}
