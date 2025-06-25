package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	dsn := "test_snippetbox?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// INFO: t.Cleanup is called when a test or subtest which uses
	// newTestDB finishes
	t.Cleanup(func() {
		script, err := os.ReadFile("testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return db
}
