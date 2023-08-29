package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"simple-go-server/handler"

	"github.com/stretchr/testify/assert"
)

func TestHandleLogin(t *testing.T) {
	assert := assert.New(t)

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlelogin1","role":"user","password":"hl1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)
	})

	t.Run("test login; success", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"handlelogin1","password":"hl1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		var login = handler.LoginResponse{}

		err := json.NewDecoder(res.Body).Decode(&login)
		assert.Nil(err)
		assert.Equal("login success", login.Message)
	})

	t.Run("test login; fail", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"handlelogin1","password":"hl4321++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusUnauthorized, res.Code)
		assert.Equal(`{"message":"wrong password"}`, res.Body.String())
	})
}

func TestHandleLogout(t *testing.T) {
	assert := assert.New(t)

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlelogout1","role":"user","password":"hl1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)
	})

	t.Run("test login; success", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"handlelogout1","password":"hl1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		var login = handler.LoginResponse{}

		err := json.NewDecoder(res.Body).Decode(&login)
		assert.Nil(err)
		assert.Equal("login success", login.Message)
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})
}
