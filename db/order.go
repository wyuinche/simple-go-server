package db

import (
	"simple-go-server/model"
	"time"

	"github.com/pkg/errors"
)

var selectOrder = `SELECT * FROM "order" WHERE oid = $1`
var insertOrder = `INSERT INTO "order" (uid, date) VALUES ($1, $2)`
var updateOrder = `UPDATE "order" SET date=$1 WHERE oid=$2`
var deleteOrder = `DELETE FROM "order" WHERE oid=$1`

var selectOrderProduct = `SELECT * FROM orderproduct WHERE oid = $1`
var insertOrderProduct = `INSERT INTO orderproduct (oid, pid) VALUES ($1, $2)`
var updateOrderProduct = `UPDATE orderproduct SET pid=$1 WHERE oid=$2 and pid=$3`
var deleteOrderProduct = `DELETE FROM orderproduct WHERE oid=$1 and pid=$2`

var selectUserOrders = `SELECT * FROM "order" WHERE uid = $1`
var selectOrders = `SELECT * FROM "order" ORDER BY date desc`

func (db *Database) InsertOrder(uid int64) (int64, error) {
	result, err := db.Exec(
		insertOrder,
		uid,
		time.Now().Unix(),
	)
	if err != nil {
		return 0, errors.Errorf("transaction execution failure")
	}

	oid, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Errorf("invalid result, no oid")
	}

	return oid, nil
}

func (db *Database) SelectOrder(oid int64) (*model.Order, error) {
	order := model.Order{}

	err := db.QueryRow(selectOrder, oid).Scan(&order.OID, &order.UID, &order.Date)
	if err == nil {
		return &order, nil
	}

	if err.Error() != "sql: no rows in result set" {
		return nil, errors.Errorf("select order failure")
	}

	return nil, nil
}

func (db *Database) UpdateOrder(oid int64) error {
	_, err := db.Exec(
		updateOrder,
		time.Now().Unix(),
		oid,
	)
	if err != nil {
		return errors.Errorf("transaction execution failure")
	}

	return nil
}

func (db *Database) DeleteOrder(oid int64) error {
	_, err := db.Exec(
		deleteOrder,
		oid,
	)
	if err != nil {
		return errors.Errorf("transaction execution failure")
	}

	return nil
}

func (db *Database) InsertOrderProduct(oid, pid int64) error {
	_, err := db.Exec(
		insertOrderProduct,
		oid,
		pid,
	)
	if err != nil {
		return errors.Errorf("transaction execution failure")
	}

	return nil
}

func (db *Database) SelectOrderProduct(oid int64) ([]model.OrderProduct, error) {
	orders := []model.OrderProduct{}

	rows, err := db.Query(selectOrderProduct, oid)
	if err != nil {
		return nil, errors.Errorf("transaction execution failure")
	}
	for {
		if !rows.Next() {
			break
		}

		order := model.OrderProduct{}
		if err = rows.Scan(&order.OID, &order.PID); err != nil {
			return nil, errors.Errorf("column scanning failure")
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (db *Database) UpdateOrderProduct(oid, oldPid, newPid int64) error {
	_, err := db.Exec(
		updateOrderProduct,
		newPid,
		oid,
		oldPid,
	)
	if err != nil {
		return errors.Errorf("transaction execution failure")
	}

	return nil
}

func (db *Database) DeleteOrderProduct(oid, pid int64) error {
	_, err := db.Exec(
		deleteOrderProduct,
		oid,
		pid,
	)
	if err != nil {
		return errors.Errorf("transaction execution failure")
	}

	return nil
}

func (db *Database) SelectUserOrders(uid int64) ([]model.Order, error) {
	orders := []model.Order{}

	rows, err := db.Query(selectUserOrders, uid)
	if err != nil {
		return nil, errors.Errorf("transaction execution failure")
	}
	for {
		if !rows.Next() {
			break
		}

		order := model.Order{}
		if err = rows.Scan(&order.OID, &order.UID, &order.Date); err != nil {
			return nil, errors.Errorf("column scanning failure")
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (db *Database) SelectOrders() ([]model.Order, error) {
	orders := []model.Order{}

	rows, err := db.Query(selectOrders)
	if err != nil {
		return nil, errors.Errorf("transaction execution failure")
	}
	for {
		if !rows.Next() {
			break
		}

		order := model.Order{}
		if err = rows.Scan(&order.OID, &order.UID, &order.Date); err != nil {
			return nil, errors.Errorf("column scanning failure")
		}

		orders = append(orders, order)
	}

	return orders, nil
}
