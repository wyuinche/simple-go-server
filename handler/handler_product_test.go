package handler_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"simple-go-server/handler"
	"simple-go-server/token"

	"github.com/stretchr/testify/assert"
)

func TestHandleCreateProduct(t *testing.T) {
	assert := assert.New(t)

	var at *http.Cookie

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlercreatep1","role":"user","password":"hcp1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		user := handler.CreateUserResponse{}

		err := json.NewDecoder(res.Body).Decode(&user)
		assert.Nil(err)
		assert.Equal("sign up success", user.Message)
	})

	t.Run("test login; manager", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"master01","password":"pwmaster01++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		var login handler.LoginResponse

		err := json.NewDecoder(res.Body).Decode(&login)
		assert.Nil(err)
		assert.Equal("login success", login.Message)

		for _, k := range res.Result().Cookies() {
			if k.Name == token.ACCESS_TOKEN_NAME {
				at = k
			}
		}
	})

	t.Run("test create product", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"cookie","price":500}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		pd := handler.CreateProductResponse{}

		err := json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})

	t.Run("test login; user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"handlercreatep1","password":"hcp1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		var login handler.LoginResponse

		err := json.NewDecoder(res.Body).Decode(&login)
		assert.Nil(err)
		assert.Equal("login success", login.Message)

		for _, k := range res.Result().Cookies() {
			if k.Name == token.ACCESS_TOKEN_NAME {
				at = k
			}
		}
	})

	t.Run("test create product; unathorized", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"cookie2","price":300}`,
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusUnauthorized, res.Code)
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})
}

func TestHandleGetProduct(t *testing.T) {
	assert := assert.New(t)
	var at *http.Cookie

	var pid int64

	t.Run("test login; manager", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"master01","password":"pwmaster01++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		var login handler.LoginResponse

		err := json.NewDecoder(res.Body).Decode(&login)
		assert.Nil(err)
		assert.Equal("login success", login.Message)

		for _, k := range res.Result().Cookies() {
			if k.Name == token.ACCESS_TOKEN_NAME {
				at = k
			}
		}
	})

	t.Run("test create product", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"ice cream","price":123}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		pd := handler.CreateProductResponse{}

		err := json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)

		pid = pd.PID
	})

	t.Run("test get product", func(t *testing.T) {
		pd := handler.GetProductResponse{}

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/product/%d", pid), nil)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		err := json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal(pid, pd.PID)
	})

	t.Run("test get product; not found", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/product/999999", nil)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusNotFound, res.Code)
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})
}

func TestHandleUpdateProduct(t *testing.T) {
	assert := assert.New(t)

	var pid int64

	var at *http.Cookie

	t.Run("test login; manager", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"master01","password":"pwmaster01++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		var login handler.LoginResponse

		err := json.NewDecoder(res.Body).Decode(&login)
		assert.Nil(err)
		assert.Equal("login success", login.Message)

		for _, k := range res.Result().Cookies() {
			if k.Name == token.ACCESS_TOKEN_NAME {
				at = k
			}
		}
	})

	t.Run("test create product", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"cookie","price":500}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		pd := handler.CreateProductResponse{}

		err := json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)

		pid = pd.PID
	})

	t.Run("test update product", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", fmt.Sprintf("/product/%d", pid), strings.NewReader(
			`{"name":"cookie22","price":1000}`,
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"update product success"}`, res.Body.String())
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})
}

func TestHandleDeleteProduct(t *testing.T) {
	assert := assert.New(t)

	var pid int64

	var at *http.Cookie

	t.Run("test login; manager", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"master01","password":"pwmaster01++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		var login handler.LoginResponse

		err := json.NewDecoder(res.Body).Decode(&login)
		assert.Nil(err)
		assert.Equal("login success", login.Message)

		for _, k := range res.Result().Cookies() {
			if k.Name == token.ACCESS_TOKEN_NAME {
				at = k
			}
		}
	})

	t.Run("test create product", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"cookie44","price":444}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		pd := handler.CreateProductResponse{}

		err := json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)

		pid = pd.PID
	})

	t.Run("test delete product", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/product/%d", pid), nil)

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"delete product success"}`, res.Body.String())
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})
}
