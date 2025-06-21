package main

// Custom type to prevent naming collisions with 3rd party packages
// that also want to use the "isAuthenticated" key
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
