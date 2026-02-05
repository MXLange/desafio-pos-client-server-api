package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// NewDBConnection establishes a new database connection using the provided connection string.
func NewDBConnection(connString string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)


	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
