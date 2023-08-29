package db

import (
	"simple-go-server/model"

	"github.com/pkg/errors"
)

var selectProduct = `SELECT * FROM product WHERE pid = $1`
var insertProduct = `INSERT INTO product (name, price) VALUES ($1, $2)`
var updateProduct = `UPDATE product SET name=$1, price=$2 WHERE pid=$3`
var deleteProduct = `DELETE FROM product WHERE pid=$1`

func (db *Database) InsertProduct(name string, price int64) (int64, error) {
	result, err := db.Exec(
		insertProduct,
		name,
		price,
	)
	if err != nil {
		return 0, errors.Errorf("transaction execution failure")
	}

	pid, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Errorf("invalid result, no pid")
	}

	return pid, nil
}

func (db *Database) SelectProduct(pid int64) (*model.Product, error) {
	product := model.Product{}

	err := db.QueryRow(selectProduct, pid).Scan(&product.PID, &product.Name, &product.Price)
	if err == nil {
		return &product, nil
	}

	if err.Error() != "sql: no rows in result set" {
		return nil, errors.Errorf("select product failure")
	}

	return nil, nil
}

func (db *Database) UpdateProduct(pid int64, name string, price int64) error {
	_, err := db.Exec(
		updateProduct,
		name,
		price,
		pid,
	)
	if err != nil {
		return errors.Errorf("transaction execution failure")
	}

	return nil
}

func (db *Database) DeleteProduct(pid int64) error {
	_, err := db.Exec(
		deleteProduct,
		pid,
	)
	if err != nil {
		return errors.Errorf("transaction execution failure")
	}

	return nil
}
