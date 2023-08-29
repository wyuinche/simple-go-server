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

func (r *Router) AddGet(k string, v gin.HandlerFunc) error {
	if _, found := r.get[k]; found {
		return APIError{
			error: errors.Errorf("such get api already exists, %s", k),
			code:  API_EC_ALREADY_EXISTS,
		}
	}

	r.get[k] = v

	return nil
}

func (r *Router) AddPost(k string, v gin.HandlerFunc) error {
	if _, found := r.post[k]; found {
		return APIError{
			error: errors.Errorf("such post api already exists, %s", k),
			code:  API_EC_ALREADY_EXISTS,
		}
	}

	r.post[k] = v

	return nil
}

func (r *Router) AddPut(k string, v gin.HandlerFunc) error {
	if _, found := r.put[k]; found {
		return APIError{
			error: errors.Errorf("such put api already exists, %s", k),
			code:  API_EC_ALREADY_EXISTS,
		}
	}

	r.put[k] = v

	return nil
}

func (r *Router) AddDelete(k string, v gin.HandlerFunc) error {
	if _, found := r.delete[k]; found {
		return APIError{
			error: errors.Errorf("such delete api already exists, %s", k),
			code:  API_EC_ALREADY_EXISTS,
		}
	}

	r.delete[k] = v

	return nil
}

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
