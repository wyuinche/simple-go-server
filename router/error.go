package router

const (
	API_EC_ALREADY_EXISTS = iota
)

type APIError struct {
	error
	code int
}

func (e APIError) Code() int {
	return e.code
}
