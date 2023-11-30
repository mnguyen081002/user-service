package domain

import (
	"context"
	"erp/internal/models"
	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (user *models.User, err error)
	Login(ctx context.Context, req LoginRequest) (res LoginResponse, err error)
}

type JwtClaims struct {
	jwt.StandardClaims
	// RoleID    string              `json:"role_id"`
	TokenType string `json:"token_type"`
}

type RegisterRequest struct {
	Email       string `json:"email" binding:"required" validate:"email"`
	Password    string `json:"password" binding:"required" validate:"min=6,max=20"`
	FirstName   string `json:"first_name" binding:"required" validate:"min=1,max=50"`
	LastName    string `json:"last_name" binding:"required" validate:"min=1,max=50"`
	RequestFrom string `json:"request_from" binding:"required" enums:"web,app"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=6,max=20"`
}

type UserInfoLoginResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginResponse struct {
	AccessToken  string                `json:"access_token"`
	RefreshToken string                `json:"refresh_token"`
	ExpiresAt    int64                 `json:"expires_at"`
	UserInfo     UserInfoLoginResponse `json:"user_info"`
}
