package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// Adds number of days to expiration
	expiration := time.Now().AddDate(0, 0, expires)
	// stmt := fmt.Sprintf("INSERT INTO snippets (title, content, created, expires) VALUES (%s, %s, DATE(), %s)",
	// 	title, content, expiration)
	stmt := "INSERT INTO snippets (title, content, created, expires) VALUES (?, ?, DATE(), ?)"
	result, err := m.DB.Exec(stmt, title, content, expiration)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Error DB.Exec: %s", err))
	}

	/* ID of our newly inserted record */
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

// This will get the most recent 10 snippets
func (m *SnippetModel) Latest() (*[]Snippet, error) {
	return nil, nil
}
