package handler

import "simple-go-server/model"

type CreateUserResponse struct {
	UID     int64  `json:"uid"`
	Message string `json:"message"`
}

type GetUserResponse struct {
	model.User
}
