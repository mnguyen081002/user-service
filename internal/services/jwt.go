package service

import (
	config "erp/config"
	"erp/internal/constants"
	"erp/internal/domain"
	"erp/internal/infrastructure"
	"erp/internal/repository"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type jwtService struct {
	ufw    *repository.UnitOfWork
	db     infrastructure.Database
	config *config.Config
	logger *zap.Logger
}

func NewJwtService(ufw *repository.UnitOfWork, db infrastructure.Database, config *config.Config, logger *zap.Logger) domain.JwtService {
	return &jwtService{
		ufw:    ufw,
		db:     db,
		config: config,
		logger: logger,
	}
}

func (j *jwtService) GenerateToken(userID string, kid string, tokenType constants.TokenType, expiresIn int64) (string, error) {
	claims := domain.JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		TokenType: string(tokenType),
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenObj.Header["kid"] = kid

	privateKeyFromFile, err := os.ReadFile("./config/keys/private.txt")
	if err != nil {
		return "", errors.WithStack(err)
	}

	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyFromFile)
	if err != nil {
		return "", errors.WithStack(err)
	}

	token, err := tokenObj.SignedString(privateKey)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}

func (j *jwtService) GenerateAuthTokens(userID string) (string, string, error) {
	accessToken, err := j.GenerateToken(userID, j.config.Jwt.Kid, constants.AccessToken, j.config.Jwt.AccessTokenExpiresIn)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.GenerateToken(userID, j.config.Jwt.Kid, constants.RefreshToken, j.config.Jwt.RefreshTokenExpiresIn)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
