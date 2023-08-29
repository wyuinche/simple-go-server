package db

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Database struct {
	*sqlx.DB
}

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
