package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

func (app *application) errorServer(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	if app.debugMode {
		http.Error(w, trace, http.StatusInternalServerError)
	} else {
		app.errorLog.Output(2, trace)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) errorClient(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) errorNotFound(w http.ResponseWriter) {
	app.errorClient(w, http.StatusNotFound)
}

func (app *application) errorDebugMode(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	w.Write([]byte(trace))
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templates[page]
	if !ok {
		err := fmt.Errorf("The template %s does not exit", page)
		app.errorServer(w, err)
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.errorServer(w, err)
		return
		/* Stop execution if there's an error with the template */
	}

	w.WriteHeader(status)

	buf.WriteTo(w)

}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.form.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		/* Invalid target destination will cause panic,
		   this is not a user error */
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		/* For all other errors we want to return and assume client error */
		return err
	}

	return nil
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
