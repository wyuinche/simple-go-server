package db

import (
	"simple-go-server/model"

	"github.com/pkg/errors"
)

var selectUserByUserID = `SELECT * FROM user WHERE userid = $1`
var insertUser = `INSERT INTO user (userid, role, password) VALUES ($1, $2, $3)`
var updateUser = `UPDATE user SET role=$1, password=$2 WHERE userid=$3`
var deleteUser = `DELETE FROM user WHERE userid=$1`

func (db *Database) SelectUser(userID string) (*model.User, error) {
	user := model.User{}

	err := db.QueryRow(selectUserByUserID, userID).Scan(&user.UID, &user.UserID, &user.Role, &user.Password)
	if err == nil {
		return &user, nil
	}

	if err.Error() != "sql: no rows in result set" {
		return nil, errors.Errorf("select user failure")
	}

	return nil, nil
}

func (db *Database) InsertUser(userID, role, pw string) (int64, error) {
	result, err := db.Exec(
		insertUser,
		userID,
		role,
		pw,
	)
	if err != nil {
		return 0, errors.Errorf("transaction execution failure")
	}

	uid, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Errorf("invalid result, no uid")
	}

	return uid, nil
}

func (db *Database) UpdateUser(userID, role, pw string) error {
	_, err := db.Exec(
		updateUser,
		role,
		pw,
		userID,
	)
	if err != nil {
		return errors.Errorf("transaction execution failure")
	}

	return nil
}

func (db *Database) DeleteUser(userID string) error {
	_, err := db.Exec(
		deleteUser,
		userID,
	)
	if err != nil {
		return errors.Errorf("transaction execution failure")
	}

	return nil
}
