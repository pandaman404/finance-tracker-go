package user

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,max=100"`
	Email string `json:"email" binding:"omitempty,email"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
