package domain

import (
	"context"
	"erp/internal/constants"
	"erp/internal/infrastructure"
	"erp/internal/models"
)

type TokenRepository interface {
	Upsert(db *infrastructure.Database, ctx context.Context, token *models.Token) (res *models.Token, err error)
}

type JwtService interface {
	GenerateToken(userID string, tokenType constants.TokenType, expiresIn int64) (string, error)
	ValidateToken(token string, tokenType constants.TokenType) (*string, error)
	GenerateAuthTokens(userID string) (string, string, error)
}
