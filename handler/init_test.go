package handler_test

import (
	"simple-go-server/handler"
	"simple-go-server/router"
)

var TestRouter *router.Router

func init() {
	r := handler.GetRouter()
	TestRouter = &r
	TestRouter.LoadAll()
}
