package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

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
	stmt := "SELECT * FROM snippets WHERE expires > DATE() AND ID = ?;"

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// This will get the most recent 10 snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := "SELECT * FROM snippets WHERE expires > DATE() ORDER BY id DESC LIMIT 10;"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	/* INFO: Optimization purposes */
	// snippets := make([]*Snippet, 10)
	snippets := []*Snippet{}

	/* INFO: The resultset will automatically close itself when iteration completes */
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	/* INFO: rows.Err() to retrieve any error that was encountered during the iteration
	   always needs to be called after rows.Next() */
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
