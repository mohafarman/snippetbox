package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	/* Error for managing user login */
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
