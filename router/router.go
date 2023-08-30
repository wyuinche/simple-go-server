package router

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Router struct {
	*gin.Engine
	get    map[string]gin.HandlerFunc
	post   map[string]gin.HandlerFunc
	put    map[string]gin.HandlerFunc
	delete map[string]gin.HandlerFunc
}

func NewRouter(e *gin.Engine) Router {
	get := make(map[string]gin.HandlerFunc)
	post := make(map[string]gin.HandlerFunc)
	put := make(map[string]gin.HandlerFunc)
	delete := make(map[string]gin.HandlerFunc)

	return Router{
		Engine: e,
		get:    get,
		post:   post,
		put:    put,
		delete: delete,
	}
}

func (r *Router) AddGet(api string, handlerFunc gin.HandlerFunc) error {
	if _, found := r.get[api]; found {
		return APIError{
			error: errors.Errorf("such get api already exists, %s", api),
			code:  API_EC_ALREADY_EXISTS,
		}
	}

	r.get[api] = handlerFunc

	return nil
}

func (r *Router) AddPost(api string, handlerFunc gin.HandlerFunc) error {
	if _, found := r.post[api]; found {
		return APIError{
			error: errors.Errorf("such post api already exists, %s", api),
			code:  API_EC_ALREADY_EXISTS,
		}
	}

	r.post[api] = handlerFunc

	return nil
}

func (r *Router) AddPut(api string, handlerFunc gin.HandlerFunc) error {
	if _, found := r.put[api]; found {
		return APIError{
			error: errors.Errorf("such put api already exists, %s", api),
			code:  API_EC_ALREADY_EXISTS,
		}
	}

	r.put[api] = handlerFunc

	return nil
}

func (r *Router) AddDelete(api string, handlerFunc gin.HandlerFunc) error {
	if _, found := r.delete[api]; found {
		return APIError{
			error: errors.Errorf("such delete api already exists, %s", api),
			code:  API_EC_ALREADY_EXISTS,
		}
	}

	r.delete[api] = handlerFunc

	return nil
}

// LoadAll sets all api handlers that r(*Router) holds (GET, POST, PUT, DELETE).
func (r *Router) LoadAll() {
	for k, v := range r.get {
		r.GET(k, v)
	}

	for k, v := range r.post {
		r.POST(k, v)
	}

	for k, v := range r.put {
		r.PUT(k, v)
	}

	for k, v := range r.delete {
		r.DELETE(k, v)
	}
}
