package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"simple-go-server/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleCreateOrder(c *gin.Context) {
	req := new(CreateOrderRequest)

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid request format")
		return
	}

	if len(req.Products) == 0 {
		writeMessage(c, http.StatusBadRequest, "empty products")
	}

	claims, keep := checkToken(c)
	if !keep {
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	for _, pid := range req.Products {
		product, err := db.SelectProduct(pid)
		if err != nil {
			writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			return
		}

		if product == nil {
			writeMessage(c, http.StatusNotFound, "product not found")
			break
		}
	}

	oid, err := db.InsertOrder(claims.UID)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	pids := []int64{}
	var insertErr error = nil
	for _, pid := range req.Products {
		err := db.InsertOrderProduct(oid, pid)
		if err != nil {
			insertErr = err
			break
		}
		pids = append(pids, pid)
	}

	if len(pids) != len(req.Products) || insertErr != nil {
		err := db.DeleteOrder(oid)
		if err != nil {
			log.Println("failed to delete order from database, oid: ", oid)
		}
		for _, pid := range pids {
			err := db.DeleteOrderProduct(oid, pid)
			if err != nil {
				log.Printf("failed to delete order product from database, oid: %d, pid: %d", oid, pid)
			}
		}
		writeMessage(c, http.StatusInternalServerError, "order product insert failure")
		return
	}

	c.JSON(
		http.StatusCreated,
		CreateOrderResponse{
			oid,
			"create order success",
		},
	)
}

func handleGetOrder(c *gin.Context) {
	oid, err := strconv.Atoi(c.Param("oid"))
	if err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid order id format")
		return
	}

	claims, keep := checkToken(c)
	if !keep {
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	order, err := db.SelectOrder(int64(oid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if order == nil {
		writeMessage(c, http.StatusNotFound, "order not found")
		return
	}

	if claims.UID != order.UID {
		writeMessage(c, http.StatusUnauthorized, "not order of user")
		return
	}

	orders, err := db.SelectOrderProduct(int64(oid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if len(orders) == 0 {
		writeMessage(c, http.StatusNotFound, "ordered product not found")
		return
	}

	products := make([]int64, len(orders))
	for i, od := range orders {
		products[i] = od.PID
	}

	c.JSON(
		http.StatusOK,
		GetOrderResponse{
			OID:      order.OID,
			UID:      order.UID,
			Products: products,
		},
	)
}

func handleUpdateOrder(c *gin.Context) {
	oid, err := strconv.Atoi(c.Param("oid"))
	if err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid order id format")
		return
	}

	req := new(UpdateOrderRequest)

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid request format")
		return
	}

	claims, keep := checkToken(c)
	if !keep {
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	order, err := db.SelectOrder(int64(oid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if order == nil {
		writeMessage(c, http.StatusNotFound, "order not found")
		return
	}

	if claims.UID != order.UID {
		writeMessage(c, http.StatusUnauthorized, "not order of user")
		return
	}

	newProducts := map[int64]struct{}{}

	for _, pid := range req.Products {
		product, err := db.SelectProduct(pid)
		if err != nil {
			writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
			return
		}

		if product == nil {
			writeMessage(c, http.StatusNotFound, "product not found")
			break
		}

		if _, found := newProducts[pid]; found {
			writeMessage(c, http.StatusBadRequest, "duplicate product found in request")
			return
		}

		newProducts[pid] = struct{}{}
	}

	orders, err := db.SelectOrderProduct(int64(oid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if len(orders) == 0 {
		writeMessage(c, http.StatusNotFound, "ordered product not found")
		return
	}

	oldProducts := map[int64]struct{}{}
	for _, od := range orders {
		oldProducts[od.PID] = struct{}{}
	}

	deletedOrderProduct := []int64{}

	var deleteErr error = nil
	for pid := range oldProducts {
		if _, found := newProducts[pid]; found {
			continue
		}

		err := db.DeleteOrderProduct(int64(oid), pid)
		if err != nil {
			deleteErr = err
			break
		}

		deletedOrderProduct = append(deletedOrderProduct, pid)
	}

	if deleteErr != nil {
		for _, pid := range deletedOrderProduct {
			err := db.InsertOrderProduct(int64(oid), pid)
			if err != nil {
				log.Printf("failed to insert order product into database, oid: %d, pid: %d", oid, pid)
			}
		}
		writeMessage(c, http.StatusInternalServerError, "order product delete failure")
		return
	}

	insertedOrderProduct := []int64{}

	var insertErr error = nil
	for pid := range newProducts {
		if _, found := oldProducts[pid]; found {
			continue
		}

		err := db.InsertOrderProduct(int64(oid), pid)
		if err != nil {
			insertErr = err
			break
		}

		insertedOrderProduct = append(insertedOrderProduct, pid)
	}

	if insertErr != nil {
		for _, pid := range insertedOrderProduct {
			err := db.DeleteOrderProduct(int64(oid), pid)
			if err != nil {
				log.Printf("failed to delete order product from database, oid: %d, pid: %d", oid, pid)
			}
		}
		writeMessage(c, http.StatusInternalServerError, "order product insert failure")
		return
	}

	err = db.UpdateOrder(int64(oid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	writeMessage(c, http.StatusOK, "update order success")
}

func handleDeleteOrder(c *gin.Context) {
	oid, err := strconv.Atoi(c.Param("oid"))
	if err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid order id format")
		return
	}

	claims, keep := checkToken(c)
	if !keep {
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	order, err := db.SelectOrder(int64(oid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if order == nil {
		writeMessage(c, http.StatusNotFound, "order not found")
		return
	}

	if claims.UID != order.UID {
		writeMessage(c, http.StatusUnauthorized, "not order of user")
		return
	}

	orders, err := db.SelectOrderProduct(int64(oid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if len(orders) == 0 {
		writeMessage(c, http.StatusNotFound, "ordered product not found")
		return
	}

	products := make([]int64, len(orders))
	for i, od := range orders {
		products[i] = od.PID
	}

	deletedOrderProduct := []int64{}

	var deleteErr error = nil
	for _, pid := range products {
		if err := db.DeleteOrderProduct(int64(oid), pid); err != nil {
			deleteErr = err
			break
		}

		deletedOrderProduct = append(deletedOrderProduct, pid)
	}

	if deleteErr != nil {
		for _, pid := range deletedOrderProduct {
			err := db.InsertOrderProduct(int64(oid), pid)
			if err != nil {
				log.Printf("failed to insert order product into database, oid: %d, pid: %d", oid, pid)
			}
		}
		writeMessage(c, http.StatusInternalServerError, "order product delete failure")
		return
	}

	err = db.DeleteOrder(int64(oid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	writeMessage(c, http.StatusOK, "delete order success")
}
