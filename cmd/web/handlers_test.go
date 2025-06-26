package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/mohafarman/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	// app.routes() returns all our real application routes, middleware and handlers
	// this is thanks to isolating all our routing in the app.routes() method
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	// app.routes() returns all our real application routes, middleware and handlers
	// this is thanks to isolating all our routing in the app.routes() method
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestUserSignup(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, body := ts.get(t, "/user/signup")

	validCSRFToken := extractCSRFToken(t, body)

	const (
		validName     = "Bob"
		validEmail    = "bob@example.com"
		validPassword = "password"
		formTag       = "<form action='/user/signup' method='POST' novalidate>"
	)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantFormTag  string
	}{
		{

			name:         "Valid submission",
			userName:     validName,
			userEmail:    validEmail,
			userPassword: validPassword,
			csrfToken:    validCSRFToken,
			wantCode:     http.StatusSeeOther,
		},
	}
	/* See the book for the rest of the examples */

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}

func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	const (
		validEmail    = "bob@example.com"
		validPassword = "password"
		formTag       = "<form action='/snippet/create' method='POST' novalidate>"
	)

	tests := []struct {
		name       string
		wantCode   int
		wantHeader string
	}{
		{
			name:       "Unauthenticated",
			wantCode:   303,
			wantHeader: "/user/login",
		},
		{
			name:     "Authenticated",
			wantCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Unauthenticated" {
				code, header, _ := ts.get(t, "/snippet/create")
				assert.Equal(t, code, tt.wantCode)
				if tt.wantHeader != "" {
					assert.StringContains(t, header.Get("Location"), tt.wantHeader)
				}
				return
			} else {
				_, _, body := ts.get(t, "/user/login")
				csrfToken := extractCSRFToken(t, body)

				form := url.Values{}
				form.Add("email", validEmail)
				form.Add("password", validPassword)
				form.Add("csrf_token", csrfToken)
				/* Form is valid */

				code, _, body := ts.postForm(t, "/user/login", form)
				t.Logf("code: %v\n body: %v\n", code, body)
				// TODO:
				/* Prints out code: 400 and Bad Request. Why? I don't know */

				code, _, body = ts.get(t, "/snippet/create")
				assert.Equal(t, code, tt.wantCode)
				assert.StringContains(t, body, formTag)
			}
		})
	}
}
