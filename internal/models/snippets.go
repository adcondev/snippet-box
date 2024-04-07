package models

import (
	"database/sql"
	"errors"
	"time"
)

// Use sql.NullType in case NULL values are allowed in DB

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	db         *sql.DB
	InsertStmt *sql.Stmt
	GetStmt    *sql.Stmt
	LatestStmt *sql.Stmt
}

func NewSnippetModel(db *sql.DB) (*SnippetModel, error) {

	insert := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	insertStmt, err := db.Prepare(insert)
	if err != nil {
		return nil, err
	}

	get := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	getStmt, err := db.Prepare(get)
	if err != nil {
		return nil, err
	}

	latest := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	latestStmt, err := db.Prepare(latest)
	if err != nil {
		return nil, err
	}

	return &SnippetModel{db, insertStmt, getStmt, latestStmt}, nil
}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	tx, err := sm.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Stmt(sm.InsertStmt).Exec(title, content, expires)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (*Snippet, error) {

	s := &Snippet{}

	err := sm.GetStmt.QueryRow(id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (sm *SnippetModel) Latest() ([]*Snippet, error) {

	rows, err := sm.LatestStmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
