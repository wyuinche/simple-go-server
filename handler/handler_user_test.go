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

func TestHandleCreateUser(t *testing.T) {
	assert := assert.New(t)

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlercreate1","role":"user","password":"hc1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		user1 := handler.CreateUserResponse{}

		err := json.NewDecoder(res.Body).Decode(&user1)
		assert.Nil(err)
		assert.Equal("sign up success", user1.Message)

		res = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlercreate2","role":"user","password":"hc1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		user2 := handler.CreateUserResponse{}

		err = json.NewDecoder(res.Body).Decode(&user2)
		assert.Nil(err)
		assert.Equal("sign up success", user2.Message)
		assert.Equal(user1.UID+1, user2.UID)
	})

	t.Run("test create user; conflict", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlercreate2","role":"user","password":"hc1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusConflict, res.Code)
	})

	t.Run("test create user; invalid role", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlercreate3","role":"customer","password":"hc1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusBadRequest, res.Code)
	})

	t.Run("test create user; invalid pw", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlercreate3","role":"user","password":"hc1234"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusBadRequest, res.Code)
	})
}

func TestHandleGetUser(t *testing.T) {
	assert := assert.New(t)

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlerget1","role":"user","password":"hg1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		user1 := handler.CreateUserResponse{}

		err := json.NewDecoder(res.Body).Decode(&user1)
		assert.Nil(err)

		res = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlerget2","role":"user","password":"hg1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		user2 := handler.CreateUserResponse{}

		err = json.NewDecoder(res.Body).Decode(&user2)
		assert.Nil(err)
	})

	t.Run("test get user", func(t *testing.T) {
		ui := handler.GetUserResponse{}

		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/user/handlerget1", nil)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		err := json.NewDecoder(res.Body).Decode(&ui)
		assert.Nil(err)
		assert.Equal("handlerget1", ui.UserID)

		res = httptest.NewRecorder()
		req = httptest.NewRequest("GET", "/user/handlerget2", nil)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		err = json.NewDecoder(res.Body).Decode(&ui)
		assert.Nil(err)
		assert.Equal("handlerget2", ui.UserID)
	})

	t.Run("test get user; not found", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/user/999999", nil)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusNotFound, res.Code)
	})
}

func TestHandleUpdateUser(t *testing.T) {
	assert := assert.New(t)

	var at *http.Cookie

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlerupdate1","role":"user","password":"hg1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		user := handler.CreateUserResponse{}

		err := json.NewDecoder(res.Body).Decode(&user)
		assert.Nil(err)
		assert.Equal("sign up success", user.Message)
	})

	t.Run("test login", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"handlerupdate1","password":"hg1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		var login = handler.LoginResponse{}

		err := json.NewDecoder(res.Body).Decode(&login)
		assert.Nil(err)
		assert.Equal("login success", login.Message)

		for _, k := range res.Result().Cookies() {
			if k.Name == token.ACCESS_TOKEN_NAME {
				at = k
			}
		}
	})

	t.Run("test update", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/user/handlerupdate1", strings.NewReader(
			`{"user_id":"handlerupdate1","role":"user","password":"hg1234++"}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"user update success"}`, res.Body.String())
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})
}

func TestHandleDeleteUser(t *testing.T) {
	assert := assert.New(t)
	var at *http.Cookie

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlerdelete1","role":"user","password":"hg1234++"}`,
		))

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		user := handler.CreateUserResponse{}

		err := json.NewDecoder(res.Body).Decode(&user)
		assert.Nil(err)
		assert.Equal("sign up success", user.Message)
	})

	t.Run("test login", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(
			`{"user_id":"handlerdelete1","password":"hg1234++"}`,
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

	t.Run("test delete", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/user/handlerdelete1", nil)
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"user delete success"}`, res.Body.String())
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})
}

func TestHandleGetUserOrder(t *testing.T) {
	assert := assert.New(t)
	var at *http.Cookie

	var pid int64
	var oid1 int64
	var oid2 int64
	var oid3 int64

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"userorders1","role":"user","password":"uo1234++"}`,
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
			`{"user_id":"userorders1","password":"uo1234++"}`,
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

	t.Run("test order", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/order", strings.NewReader(
			fmt.Sprintf(`{"products":[%d]}`, pid),
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		var or handler.CreateOrderResponse

		err := json.NewDecoder(res.Body).Decode(&or)
		assert.Nil(err)
		assert.Equal("create order success", or.Message)

		oid1 = or.OID

		res = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/order", strings.NewReader(
			fmt.Sprintf(`{"products":[%d]}`, pid),
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		err = json.NewDecoder(res.Body).Decode(&or)
		assert.Nil(err)
		assert.Equal("create order success", or.Message)

		oid2 = or.OID

		res = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/order", strings.NewReader(
			fmt.Sprintf(`{"products":[%d]}`, pid),
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		err = json.NewDecoder(res.Body).Decode(&or)
		assert.Nil(err)
		assert.Equal("create order success", or.Message)

		oid3 = or.OID
	})

	t.Run("test get user orders", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/user/userorders1/orders", nil)

		req.AddCookie(at)

		var os handler.GetUserOrdersResponse

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		err := json.NewDecoder(res.Body).Decode(&os)
		assert.Nil(err)
		assert.Equal(3, len(os.Orders))

		found := map[string]struct{}{}
		for _, o := range os.Orders {
			s := strings.Split(o, ",")
			assert.Len(s, 2)
			found[s[0]] = struct{}{}
		}

		_, ok := found[fmt.Sprintf("%d", oid1)]
		assert.True(ok)
		_, ok = found[fmt.Sprintf("%d", oid2)]
		assert.True(ok)
		_, ok = found[fmt.Sprintf("%d", oid3)]
		assert.True(ok)
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})
}
