package middlewares

import (
	"context"
	"erp/internal/api/response"
	"erp/internal/api_errors"
	"erp/internal/domain"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (e *GinMiddleware) Auth(authorization bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")

		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ResponseError{
				Message: "Unauthorized",
				Code:    api_errors.ErrUnauthorizedAccess,
			})
			return
		}
		jwtToken := strings.Split(auth, " ")[1]

		if jwtToken == "" {
			c.Errors = append(c.Errors, &gin.Error{
				Err: errors.New(api_errors.ErrTokenMissing),
			})

			mas := api_errors.MapErrorCodeMessage[api_errors.ErrTokenMissing]
			c.AbortWithStatusJSON(mas.Status, response.ResponseError{
				Message: mas.Message,
				Code:    api_errors.ErrTokenMissing,
			})
			return
		}

		claims, err := parseToken(jwtToken, e.config.Jwt.Secret)
		if err != nil {
			c.Errors = append(c.Errors, &gin.Error{
				Err: errors.WithStack(err),
			})
			mas := api_errors.MapErrorCodeMessage[err.Error()]
			c.AbortWithStatusJSON(mas.Status, response.ResponseError{
				Message: mas.Message,
				Code:    api_errors.ErrTokenInvalid,
			})
			return
		}

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "x-user-id", claims.Subject))
		if !authorization {
			c.Next()
			return
		}

		storeID := c.Request.Header.Get("x-store-id")
		if storeID == "" {
			c.Errors = append(c.Errors, &gin.Error{
				Err: errors.New(api_errors.ErrMissingXStoreID),
			})

			mas := api_errors.MapErrorCodeMessage[api_errors.ErrMissingXStoreID]

			c.AbortWithStatusJSON(mas.Status, response.ResponseError{
				Message: mas.Message,
				Code:    api_errors.ErrMissingXStoreID,
			})
			return
		}
		c.Next()
	}
}

func parseToken(jwtToken string, secret string) (*domain.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &domain.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		if (err.(*jwt.ValidationError)).Errors == jwt.ValidationErrorExpired {
			return nil, errors.New(api_errors.ErrTokenExpired)
		}
		return nil, errors.Wrap(err, "cannot parse token")
	}

	if claims, OK := token.Claims.(*domain.JwtClaims); OK && token.Valid {
		return claims, nil
	}

	return nil, errors.New(api_errors.ErrTokenInvalid)
}
