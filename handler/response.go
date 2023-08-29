package handler

import "simple-go-server/model"

type LoginResponse struct {
	UID     int64  `json:"uid"`
	Message string `json:"message"`
}

type CreateUserResponse struct {
	UID     int64  `json:"uid"`
	Message string `json:"message"`
}

type GetUserResponse struct {
	model.User
}

type CreateProductResponse struct {
	PID     int64  `json:"pid"`
	Message string `json:"message"`
}

type GetProductResponse struct {
	model.Product
}

type CreateOrderResponse struct {
	OID     int64  `json:"oid"`
	Message string `json:"message"`
}

type GetOrderResponse struct {
	OID      int64   `json:"oid"`
	UID      int64   `json:"uid"`
	Products []int64 `json:"products"`
	Date     string  `json:"date"`
}
