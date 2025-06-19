package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3" /* INFO: Uses CGO */
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, DATE())"

	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		/* Handle duplicate email adress */
		var sqlite3Error sqlite3.Error
		if errors.As(err, &sqlite3Error) {
			if errors.Is(sqlite3Error.Code, sqlite3.ErrConstraint) {
				return ErrDuplicateEmail
			}
		}
		/* Simply return err for all other errors */
		return err
	}
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
