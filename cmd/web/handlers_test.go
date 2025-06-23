package main

import (
	"net/http"
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
