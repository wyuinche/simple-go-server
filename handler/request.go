package handler

type CreateUserRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	Password string `json:"password"`
	Role     string `json:"role"`
}
