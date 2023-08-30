package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	*sqlx.DB
}

// Connect connects the database or verifies
// that the database is connected normally by sqlx.DB.Ping().
func (db *Database) Connect() error {
	if db.DB == nil {
		d, err := sqlx.Connect("sqlite3", ":memory:")
		if err != nil {
			return err
		}
		db.DB = d
	}

	return db.Ping()
}

// Exec executes the transaction for the queries it receives
// and returns the result as sql.Result.
func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	res, err := tx.Exec(
		query,
		args...,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return res, nil
}
