package controller

import (
	"erp/internal/domain"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

type UserController struct {
	userService domain.UserService
	logger      *zap.Logger
}

func NewUserController(userService domain.UserService, logger *zap.Logger) *UserController {
	controller := &UserController{
		userService: userService,
		logger:      logger,
	}
	return controller
}

func (b *UserController) ListUsers(c *gin.Context) {
	var req domain.ListUsersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		ResponseValidationError(c, err)
		return
	}

	res, total, err := b.userService.ListUsers(req.PageOptions, c.Request.Context())

	if err != nil {
		ResponseError(c, err)
		return
	}
	ResponseList(c, "success", total, res)
}
