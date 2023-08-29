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

func TestHandleCreateOrder(t *testing.T) {
	assert := assert.New(t)

	var at *http.Cookie
	pids := []int64{}

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlercreateo1","role":"user","password":"hco1234++"}`,
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
			`{"name":"cookie1","price":500}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		pd := handler.CreateProductResponse{}

		err := json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)

		pids = append(pids, pd.PID)

		res = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"cookie2","price":1000}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		pd = handler.CreateProductResponse{}

		err = json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)

		pids = append(pids, pd.PID)
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
			`{"user_id":"handlercreateo1","password":"hco1234++"}`,
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
			fmt.Sprintf(`{"products":[%d,%d]}`, pids[0], pids[1]),
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		var or handler.CreateOrderResponse

		err := json.NewDecoder(res.Body).Decode(&or)
		assert.Nil(err)
		assert.Equal("create order success", or.Message)
	})

	t.Run("test order; product not found", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/order", strings.NewReader(
			`{"products":[9999999999344,9999999234999]}`,
		))

		req.AddCookie(at)

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

func TestHandleGetOrder(t *testing.T) {
	assert := assert.New(t)
	var at *http.Cookie

	pids := []int64{}
	var oid int64
	var uid int64

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlergeto1","role":"user","password":"hgo1234++"}`,
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

		pids = append(pids, pd.PID)

		res = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"ice cream2","price":123}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		err = json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)

		pids = append(pids, pd.PID)
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
			`{"user_id":"handlergeto1","password":"hgo1234++"}`,
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

		uid = login.UID
	})

	t.Run("test order", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/order", strings.NewReader(
			fmt.Sprintf(`{"products":[%d,%d]}`, pids[0], pids[1]),
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		var or handler.CreateOrderResponse

		err := json.NewDecoder(res.Body).Decode(&or)
		assert.Nil(err)
		assert.Equal("create order success", or.Message)

		oid = or.OID
	})

	t.Run("test get order", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/order/%d", oid), nil)

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		var or handler.GetOrderResponse

		err := json.NewDecoder(res.Body).Decode(&or)
		assert.Nil(err)
		assert.Equal(uid, or.UID)
		assert.Equal(oid, or.OID)
		assert.Len(or.Products, 2)
	})

	t.Run("test get order; not found", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/order/999999", nil)

		req.AddCookie(at)

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

func TestHandleUpdateOrder(t *testing.T) {
	assert := assert.New(t)
	var at *http.Cookie

	var pid1 int64
	var pid2 int64
	var pid3 int64
	var oid int64

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlerupdateo1","role":"user","password":"huo1234++"}`,
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

		pid1 = pd.PID

		res = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"ice cream2","price":123}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		pd = handler.CreateProductResponse{}

		err = json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)

		pid2 = pd.PID

		res = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"ice cream3","price":123}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		pd = handler.CreateProductResponse{}

		err = json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)

		pid3 = pd.PID
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
			`{"user_id":"handlerupdateo1","password":"huo1234++"}`,
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
			fmt.Sprintf(`{"products":[%d,%d]}`, pid1, pid2),
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		var or handler.CreateOrderResponse

		err := json.NewDecoder(res.Body).Decode(&or)
		assert.Nil(err)
		assert.Equal("create order success", or.Message)

		oid = or.OID
	})

	t.Run("test update order", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", fmt.Sprintf("/order/%d", oid), strings.NewReader(
			fmt.Sprintf(`{"products":[%d,%d]}`, pid2, pid3),
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"update order success"}`, res.Body.String())
	})

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
	})
}

func TestHandleDeleteOrder(t *testing.T) {
	assert := assert.New(t)
	var at *http.Cookie

	var pid1 int64
	var pid2 int64
	var oid int64

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlerdeleteo1","role":"user","password":"hdo1234++"}`,
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

		pid1 = pd.PID

		res = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "/product", strings.NewReader(
			`{"name":"ice cream2","price":123}`,
		))
		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		pd = handler.CreateProductResponse{}

		err = json.NewDecoder(res.Body).Decode(&pd)
		assert.Nil(err)
		assert.Equal("register product success", pd.Message)

		pid2 = pd.PID
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
			`{"user_id":"handlerdeleteo1","password":"hdo1234++"}`,
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
			fmt.Sprintf(`{"products":[%d,%d]}`, pid1, pid2),
		))

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusCreated, res.Code)

		var or handler.CreateOrderResponse

		err := json.NewDecoder(res.Body).Decode(&or)
		assert.Nil(err)
		assert.Equal("create order success", or.Message)

		oid = or.OID
	})

	t.Run("test delete order", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/order/%d", oid), nil)

		req.AddCookie(at)

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"delete order success"}`, res.Body.String())
	})

	t.Run("test get order; not found", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/order/%d", oid), nil)

		req.AddCookie(at)

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

func TestHandleGetOrders(t *testing.T) {
	assert := assert.New(t)
	var at *http.Cookie

	var pid int64
	var oid1 int64
	var oid2 int64
	var oid3 int64

	t.Run("test create user", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/user", strings.NewReader(
			`{"user_id":"handlerorders1","role":"user","password":"ho1234++"}`,
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
			`{"user_id":"handlerorders1","password":"ho1234++"}`,
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

	t.Run("test logout", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/logout", nil)

		TestRouter.ServeHTTP(res, req)

		assert.Equal(http.StatusOK, res.Code)
		assert.Equal(`{"message":"logout success"}`, res.Body.String())
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

	t.Run("test get orders", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/orders", nil)

		req.AddCookie(at)

		var os handler.GetOrdersResponse

		TestRouter.ServeHTTP(res, req)
		assert.Equal(http.StatusOK, res.Code)

		err := json.NewDecoder(res.Body).Decode(&os)
		assert.Nil(err)
		assert.True(len(os.Orders) > 2)

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
