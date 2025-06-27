package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/mattn/go-sqlite3" /* INFO: Uses CGO */
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*User, error)
	CompareAndUpdatePassword(id int, currentPassword, newPassword string) (bool, error)
}

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
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		/* Error thrown when there are no rows:
		   "sql: no rows in result set" */
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, nil
		}
	}

	return id, nil
}

/* Checks if user already exists in the table 'users' */
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}

func (m *UserModel) Get(id int) (*User, error) {
	user := &User{}

	stmt := "SELECT id, name, email, created FROM users WHERE id = ?;"

	row := m.DB.QueryRow(stmt, id)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (m *UserModel) CompareAndUpdatePassword(id int, currentPassword, newPassword string) (bool, error) {
	var hashedDBpassword string

	stmt := "SELECT hashed_password FROM users WHERE id = ?;"

	row := m.DB.QueryRow(stmt, id)
	err := row.Scan(&hashedDBpassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedDBpassword), []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, ErrInvalidCredentials
		}
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return false, err
	}

	stmt = "UPDATE users SET hashed_password = ?"
	m.DB.Exec(stmt, hashedNewPassword)
	if err != nil {
		return false, err
	}

	return true, nil
}
