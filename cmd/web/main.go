package main

import (
	"log/slog"
	"net/http"
	"time"
)

var (
	PORT = ":4000"
)

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
