package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"simple-go-server/db"
	"simple-go-server/model"
	"simple-go-server/token"

	"github.com/gin-gonic/gin"
)

func handleLogin(c *gin.Context) {
	req := new(LoginRequest)

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid request format")
		return
	}

	userID := model.UserID(req.UserID)
	if err := userID.IsValid(); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid user id format")
		return
	}

	pw := model.Password(req.Password)
	if err := pw.IsValid(); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid password format")
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	user, err := db.SelectUser(string(userID))
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if user == nil {
		writeMessage(c, http.StatusNotFound, "user not found")
		return
	}

	verified := model.Password(req.Password).CompareWithHash(string(user.Password))
	if !verified {
		writeMessage(c, http.StatusUnauthorized, "wrong password")
		return
	}

	at, err := token.CreateAccessToken(user.UID, user.UserID, user.Role)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "jwt failure")
		return
	}

	c.SetCookie(token.ACCESS_TOKEN_NAME, at, 3600, "/", "localhost", false, true)

	c.JSON(
		http.StatusOK,
		LoginResponse{
			user.UID,
			"login success",
		},
	)
}

func handleLogout(c *gin.Context) {
	c.SetCookie(token.ACCESS_TOKEN_NAME, "", -1, "/", "localhost", false, true)
	writeMessage(c, http.StatusOK, "logout success")
}
