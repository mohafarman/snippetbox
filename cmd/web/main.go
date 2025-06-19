package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mohafarman/snippetbox/internal/models"
)

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templates      map[string]*template.Template
	form           *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	addr := flag.String("port", "4000", "HTTP server port adress.")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	/* INFO: driver-specific parameter which instructs
	our driver to convert SQL TIME and DATE fields to Go time.Time objects */
	dsn := "snippetbox?parseTime=true"
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templates, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionsManager := scs.New()
	sessionsManager.Store = sqlite3store.New(db)
	sessionsManager.Lifetime = 12 * time.Hour
	/* Cookie will only be sent by a users browser when there is an HTTPS connection */
	sessionsManager.Cookie.Secure = true

	app := &application{
		infoLog:  infoLog,
		errorLog: errorLog,
		snippets: &models.SnippetModel{
			DB: db,
		},
		users: &models.UserModel{
			DB: db,
		},
		templates:      templates,
		form:           formDecoder,
		sessionManager: sessionsManager,
	}

	/* Curve preferences value, so that only elliptic curves with
	   assembly implementations are used. This is because the others (as of Go 1.20)
	   CPU intensive. If we omit them the server will be more performant under
	   heavy load. */
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := http.Server{
		Addr:         ":" + *addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
	}

	infoLog.Print("Listening on port " + *addr)
	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Print("Server failed to run. Error: " + err.Error())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
