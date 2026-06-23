package dto

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
