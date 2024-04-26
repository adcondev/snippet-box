package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
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
	DB         *sql.DB
	InsertStmt *sql.Stmt
	AuthStmt   *sql.Stmt
	ExistsStmt *sql.Stmt
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

func NewUserModel(db *sql.DB) (*UserModel, error) {

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

	return &UserModel{db, insertStmt, authStmt, existsStmt}, nil
}

func (um *UserModel) Insert(name, email, password string) error {

	// Start a new transaction.
	// If there's an error (for example, if the transaction can't be started), return 0 and the error.
	tx, err := um.DB.Begin()
	if err != nil {
		return err
	}

	// Use the defer keyword to ensure that the transaction is rolled back if any subsequent code returns an error.
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = um.InsertStmt.Exec(name, email, hashedPassword)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	// Commit the transaction.
	// If there's an error (for example, if the transaction can't be committed), return 0 and the error.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (um *UserModel) Authenticate(email, password string) (int, error) {

	var id int
	var hashedPassword []byte

	err := um.AuthStmt.QueryRow(email).Scan(&id, &hashedPassword)
	if err != nil {
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
			return 0, err
		}
	}

	return id, nil
}

func (um *UserModel) Exists(id int) (bool, error) {

	var exists bool

	err := um.ExistsStmt.QueryRow(id).Scan(&exists)

	return exists, err
}
