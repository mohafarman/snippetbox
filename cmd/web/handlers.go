package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mohafarman/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/* Serve 404 not found if it's not root */
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.errorServer(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// INFO: To extract from the url, new http module in Go 1.22 allows params
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		app.errorNotFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.errorNotFound(w)
		} else {
			app.errorServer(w, err)
		}
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	/* Only allow POST method */
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// INFO: Suppress header, the w.Header().Del("Date") does not remove header.
		// w.Header()["Date"] = nil
		app.errorClient(w, http.StatusMethodNotAllowed)
		// http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "Mr snail"
	content := "Mr snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.errorServer(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
