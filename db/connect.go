package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

var db *Database

var createUserTableQuery = `CREATE TABLE user (
	uid integer primary key autoincrement,
	name text,
	email text,
	role text,
	password text);`
var createProductTableQuery = `CREATE TABLE product (
	pid integer primary key autoincrement,
	name text,
	price integer,
	stock integer);`
var createOrderTableQuery = `CREATE TABLE "order" (
	oid integer primary key autoincrement,
	uid integer,
	status text,
	date integer);`

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

func Get() (*Database, error) {
	if err := Init(); err != nil {
		return nil, err
	}

	if db == nil {
		return nil, errors.New("db disconnected")
	}

	return db, nil
}

func Init() error {
	if db != nil {
		return nil
	}

	db = new(Database)

	if err := db.Connect(); err != nil {
		return err
	}

	if _, err := db.Exec(createUserTableQuery); err != nil {
		return err
	}

	if _, err := db.Exec(createProductTableQuery); err != nil {
		return err
	}

	if _, err := db.Exec(createOrderTableQuery); err != nil {
		return err
	}

	return nil
}
