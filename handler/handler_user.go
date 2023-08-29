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

func handleCreateUser(c *gin.Context) {
	req := new(CreateUserRequest)

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

	if req.Role != model.RoleUser && req.Role != model.RoleManager {
		writeMessage(c, http.StatusBadRequest, "invalid role (not 'user' neither 'menager')")
	}

	if req.Role == model.RoleManager {
		claims, keep := checkToken(c)
		if !keep {
			return
		}

		if claims.Role != model.RoleManager {
			writeMessage(c, http.StatusUnauthorized, "general user cannot create manager account")
			return
		}
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	user, err := db.SelectUser(req.UserID)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if user != nil {
		writeMessage(c, http.StatusConflict, "already registered user")
		return
	}

	pwHash, err := pw.Hash()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "password hashing failure")
		return
	}

	uid, err := db.InsertUser(req.UserID, req.Role, pwHash)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	c.JSON(
		http.StatusCreated,
		CreateUserResponse{
			uid,
			"sign up success",
		},
	)
}

func handleGetUser(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		writeMessage(c, http.StatusBadRequest, "empty user id")
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	user, err := db.SelectUser(userID)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if user == nil {
		writeMessage(c, http.StatusNotFound, "user not found")
		return
	}

	c.JSON(
		http.StatusOK,
		GetUserResponse{*user},
	)
}

func handleUpdateUser(c *gin.Context) {
	userID := c.Param("user_id")

	req := new(UpdateUserRequest)
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid request format")
		return
	}

	claims, keep := checkToken(c)
	if !keep {
		return
	}

	if claims.UserID != userID {
		writeMessage(c, http.StatusUnauthorized, "invalid access token for this user")
		return
	}

	if req.Role != claims.Role {
		if req.Role == model.RoleManager && claims.Role != model.RoleManager {
			writeMessage(c, http.StatusUnauthorized, "user cannot become manager itself")
			return
		}
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	user, err := db.SelectUser(userID)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if user == nil {
		writeMessage(c, http.StatusNotFound, "user not found")
		return
	}

	if req.Role != user.Role {
		if req.Role == model.RoleManager && user.Role != model.RoleManager {
			writeMessage(c, http.StatusUnauthorized, "user cannot become manager itself")
			return
		}
	}

	pw := model.Password(req.Password)
	if err := pw.IsValid(); err != nil {
		writeMessage(c, http.StatusBadRequest, "invalid password format")
		return
	}

	pwHash, err := pw.Hash()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "password hashing failure")
		return
	}

	err = db.UpdateUser(userID, req.Role, pwHash)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	writeMessage(c, http.StatusOK, "user update success")
}

func handleDeleteUser(c *gin.Context) {
	userID := c.Param("user_id")

	claims, keep := checkToken(c)
	if !keep {
		return
	}

	if claims.UserID != userID {
		writeMessage(c, http.StatusUnauthorized, "invalid access token for this user")
		return
	}

	db, err := db.Get()
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, "db failure")
		return
	}

	user, err := db.SelectUser(userID)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	if user == nil {
		writeMessage(c, http.StatusNotFound, "user not found")
		return
	}

	err = db.DeleteUser(userID)
	if err != nil {
		writeMessage(c, http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	c.SetCookie(token.ACCESS_TOKEN_NAME, "", -1, "/", "localhost", false, true)

	writeMessage(c, http.StatusOK, "user delete success")
}
