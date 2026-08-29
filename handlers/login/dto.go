package login

import "github.com/golang-jwt/jwt/v4"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomClaims struct {
	Role   string `json:"role"`
	UserID int64  `json:"user_id"`
	jwt.RegisteredClaims
}

type SignupRequest struct {
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Role         string `json:"role"`
	MobileNumber string `json:"mobile_num"`
	Password     string `json:"password"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
}
