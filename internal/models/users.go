package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

/* Return user ID */
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

/* Checks if user already exists in the table 'users' */
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
