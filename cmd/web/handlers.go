package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/mohafarman/snippetbox/internal/models"
	"github.com/mohafarman/snippetbox/internal/validator"
)

type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userChangePasswordForm struct {
	Current_Password     string `form:"current_password"`
	New_Password         string `form:"new_password"`
	Confirm_New_Password string `form:"confirm_new_password"`
	validator.Validator  `form:"-"`
}

func (app *application) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.errorServer(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "about.tmpl.html", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// INFO: To extract from the url, new http module in Go 1.22 allows params
	// like so: id, err := strconv.Atoi(r.PathValue("id"))
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
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
	data := app.newTemplateData(r)
	/* INFO: data.Form has to be initialized or it is nil and will cause a
	   500 internal server error. Also good time to set default values */
	data.Form = SnippetCreateForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form SnippetCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.errorClient(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters.")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365.")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.errorServer(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet succesfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}

	app.render(w, http.StatusOK, "signup.tmpl.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.errorClient(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email adress")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl.html", data)
			return
		} else {
			app.errorServer(w, err)
		}

		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please sign in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var id int
	form := userLoginForm{}

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.errorClient(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email adress")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err = app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.errorServer(w, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.errorServer(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	redirect := app.sessionManager.PopString(r.Context(), "redirect")
	if redirect != "" {
		http.Redirect(w, r, redirect, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	user, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
	}

	data.User = user

	app.render(w, http.StatusOK, "account.tmpl.html", data)
}

func (app *application) changePasswordView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userChangePasswordForm{}
	app.render(w, http.StatusOK, "changepassword.tmpl.html", data)
}

func (app *application) changePasswordPost(w http.ResponseWriter, r *http.Request) {
	var form userChangePasswordForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.errorClient(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Current_Password), "current_password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Current_Password, 8), "current_password", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.New_Password), "new_password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.New_Password, 8), "new_password", "This field must be at least 8 characters long")
	form.CheckField(validator.NotBlank(form.Confirm_New_Password), "confirm_new_password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Confirm_New_Password, 8), "confirm_new_password", "This field must be at least 8 characters long")
	form.CheckField(form.New_Password == form.Confirm_New_Password, "confirm_new_password", "Passwords do not match")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "changepassword.tmpl.html", data)
		return
	}

	id := app.sessionManager.Get(r.Context(), "authenticatedUserID").(int)
	/* Compare passwords and Update db */
	if ok, err := app.users.CompareAndUpdatePassword(id, form.Current_Password, form.New_Password); !ok {
		if errors.Is(err, models.ErrInvalidCredentials) {
			data := app.newTemplateData(r)
			form.AddFieldError("current_password", "Current password is incorrect")
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "changepassword.tmpl.html", data)
			return
		} else {
			app.errorServer(w, err)
			return
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Password changed successfully!")

	http.Redirect(w, r, "/user/account", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.errorServer(w, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been successfully logged out")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		// User does not exist in the db
		return false
	}

	// Should be true
	return isAuthenticated
}
