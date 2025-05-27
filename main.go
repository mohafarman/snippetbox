package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

var (
	PORT = ":4000"
)

func home(w http.ResponseWriter, r *http.Request) {
	/* Serve 404 not found if it's not root */
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from Snippetbox!"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	// INFO: To extract from the url, new http module in Go 1.22 allows params
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d\n", id)
}
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	/* Only allow POST method */
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// INFO: Suppress header, the w.Header().Del("Date") does not remove header.
		// w.Header()["Date"] = nil
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a specific snippet..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view/{id}", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	server := http.Server{
		Addr:         PORT,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  90 * time.Second,
		Handler:      mux,
	}

	slog.Info("Listening on port " + PORT)
	err := server.ListenAndServe()
	slog.Error("Server failed to run. Error: " + err.Error())
}
