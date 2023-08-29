package handler

type LoginRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateProductRequest struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type UpdateProductRequest struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}
