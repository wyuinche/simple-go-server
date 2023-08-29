package handler

import (
	"net/http"
	"simple-go-server/router"
	"simple-go-server/token"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetRouter() router.Router {
	r := router.NewRouter(gin.Default())

	r.NoRoute(func(c *gin.Context) {
		writeMessage(c, http.StatusNotFound, "page not found")
	})

	r.AddGet("/", handlePing)

	r.AddPost("/user", handleCreateUser) // user - post

	r.AddGet("/user/:user_id", handleGetUser)
	r.AddPut("/user/:user_id", handleUpdateUser)
	r.AddDelete("/user/:user_id", handleDeleteUser)

	r.AddGet("/user/:user_id/orders", handleGetUserOrders)

	r.AddPost("/login", handleLogin)
	r.AddPost("/logout", handleLogout)

	r.AddPost("/product", handleCreateProduct)

	r.AddGet("/product/:pid", handleGetProduct)
	r.AddPut("/product/:pid", handleUpdateProduct)
	r.AddDelete("/product/:pid", handleDeleteProduct)

	r.AddPost("/order", handleCreateOrder)

	r.AddGet("/order/:oid", handleGetOrder)
	r.AddPut("/order/:oid", handleUpdateOrder)
	r.AddDelete("/order/:oid", handleDeleteOrder)

	r.AddGet("/orders", handleGetOrders)

	return r
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func checkToken(c *gin.Context) (*token.Claims, bool) {
	accessToken, err := c.Cookie(token.ACCESS_TOKEN_NAME)
	if err != nil {
		if err == http.ErrNoCookie {
			writeMessage(c, http.StatusUnauthorized, "no cookie")
			return nil, false
		}
		writeMessage(c, http.StatusInternalServerError, "lookup cookie failure")
		return nil, false
	}

	claims, t, err := token.GetJWTToken(accessToken)
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			writeMessage(c, http.StatusUnauthorized, "invalid jwt signature")
			return nil, false
		}
		writeMessage(c, http.StatusInternalServerError, "jwt signature check failure")
		return nil, false
	}

	if !t.Valid {
		writeMessage(c, http.StatusUnauthorized, "invalid jwt")
		return nil, false
	}

	return claims, true
}

func writeMessage(c *gin.Context, code int, msg string) {
	c.JSON(
		code,
		gin.H{
			"message": msg,
		},
	)
}
