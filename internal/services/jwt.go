package service

import (
	config "erp/config"
	"erp/internal/api_errors"
	"erp/internal/constants"
	"erp/internal/domain"
	"erp/internal/infrastructure"
	"erp/internal/repository"
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

func (j *jwtService) GenerateToken(userID string, tokenType constants.TokenType, expiresIn int64) (string, error) {
	claims := domain.JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   userID,
			ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		TokenType: string(tokenType),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.config.Jwt.Secret))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return token, nil
}

func (j *jwtService) GenerateAuthTokens(userID string) (string, string, error) {

	j.logger.Debug("Generating auth tokens", zap.Any("ExpiresIn", j.config.Jwt.AccessTokenExpiresIn))

	accessToken, err := j.GenerateToken(userID, constants.AccessToken, j.config.Jwt.AccessTokenExpiresIn)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.GenerateToken(userID, constants.RefreshToken, j.config.Jwt.RefreshTokenExpiresIn)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (j *jwtService) ValidateToken(token string, tokenType constants.TokenType) (*string, error) {
	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.Jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New(api_errors.ErrTokenExpired)
	}

	if claims.Issuer != "erp" {
		return nil, errors.New(api_errors.ErrTokenInvalid)
	}

	if claims.Subject == "" {
		return nil, errors.New(api_errors.ErrTokenInvalid)
	}

	return &claims.Subject, nil
}
