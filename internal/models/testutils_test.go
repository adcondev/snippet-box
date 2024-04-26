package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestUserModel(t *testing.T) (*UserModel, error) {

	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	insert := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	insertStmt, err := db.Prepare(insert)
	if err != nil {
		return nil, err
	}

	auth := `SELECT id, hashed_password FROM users WHERE email = ?`

	authStmt, err := db.Prepare(auth)
	if err != nil {
		return nil, err
	}

	exists := `SELECT EXISTS(SELECT true FROM users WHERE id = ?)`

	existsStmt, err := db.Prepare(exists)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {

		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	return &UserModel{db, insertStmt, authStmt, existsStmt}, nil
}
