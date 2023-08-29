package db

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var db *Database

var createUserTableQuery = `CREATE TABLE user (
	uid integer primary key autoincrement,
	userid text,
	role text,
	password text);`
var createProductTableQuery = `CREATE TABLE product (
	pid integer primary key autoincrement,
	name text,
	price integer);`
var createOrderTableQuery = `CREATE TABLE "order" (
	oid integer primary key autoincrement,
	uid integer,
	date integer);`
var createOrderProductQuery = `CREATE TABLE orderproduct (
	oid integer,
	pid integer);`

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

	if _, err := db.Exec(createOrderProductQuery); err != nil {
		return err
	}

	return nil
}
