package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	/* Serve 404 not found if it's not root */
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.errorServer(w, err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorServer(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// INFO: To extract from the url, new http module in Go 1.22 allows params
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.errorNotFound(w)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d\n", id)
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
	w.Write([]byte("Create a specific snippet..."))
}
