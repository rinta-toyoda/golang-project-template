package dto

type SignupRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Phone    string `json:"phone"    binding:"required,e164"`
	Password string `json:"password" binding:"required,min=8"`
}

type SignupResponse struct {
	Message  string `json:"message"`
	JWTToken string `json:"jwt_token"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message  string `json:"message"`
	JWTToken string `json:"jwt_token"`
}
