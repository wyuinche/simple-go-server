package handler_test

import (
	"simple-go-server/handler"
	"simple-go-server/router"
)

var TestRouter *router.Router

// init initiate the router used to test handlers
// before the test starts.
func init() {
	r := handler.GetRouter()
	TestRouter = &r
	TestRouter.LoadAll()
}
