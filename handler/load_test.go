package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"simple-go-server/db"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	assert := assert.New(t)
	err := db.Init()
	assert.Nil(err)
}

func TestHandlePing(t *testing.T) {
	expectedResponse := `{"message":"pong"}`

	assert := assert.New(t)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	TestRouter.ServeHTTP(res, req)
	assert.Equal(http.StatusOK, res.Code)

	response, err := io.ReadAll(res.Body)
	assert.Nil(err)
	assert.Equal(expectedResponse, string(response))
}
