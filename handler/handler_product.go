package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple-go-server/db"
	"simple-go-server/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleCreateProduct(c *gin.Context) {
	req := new(CreateProductRequest)

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid request format")
		return
	}

	name := model.ProductName(req.Name)
	if err := name.IsValid(); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid product name format")
		return
	}

	claims, keep := checkToken(c)
	if !keep {
		return
	}

	if claims.Role != model.RoleManager {
		writeMessage(c, http.StatusUnauthorized, "general user cannot register product")
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	pid, err := db.InsertProduct(req.Name, req.Price)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	c.JSON(
		http.StatusCreated,
		CreateProductResponse{
			pid,
			"register product success",
		},
	)
}

func handleGetProduct(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid product id format")
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	product, err := db.SelectProduct(int64(pid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if product == nil {
		writeMessage(c, http.StatusNotFound, "product not found")
		return
	}

	c.JSON(
		http.StatusOK,
		GetProductResponse{*product},
	)
}

func handleUpdateProduct(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid product id format")
		return
	}

	req := new(UpdateProductRequest)

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid request format")
		return
	}

	name := model.ProductName(req.Name)
	if err := name.IsValid(); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid product name format")
		return
	}

	claims, keep := checkToken(c)
	if !keep {
		return
	}

	if claims.Role != model.RoleManager {
		writeMessage(c, http.StatusUnauthorized, "general user cannot update product")
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	product, err := db.SelectProduct(int64(pid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if product == nil {
		writeMessage(c, http.StatusNotFound, "product not found")
		return
	}

	err = db.UpdateProduct(int64(pid), req.Name, req.Price)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	writeMessage(c, http.StatusOK, "update product success")
}

func handleDeleteProduct(c *gin.Context) {
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid product id format")
		return
	}

	claims, keep := checkToken(c)
	if !keep {
		return
	}

	if claims.Role != model.RoleManager {
		writeMessage(c, http.StatusUnauthorized, "general user cannot delete product")
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	product, err := db.SelectProduct(int64(pid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if product == nil {
		writeMessage(c, http.StatusNotFound, "product not found")
		return
	}

	err = db.DeleteProduct(int64(pid))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	writeMessage(c, http.StatusOK, "delete product success")
}
