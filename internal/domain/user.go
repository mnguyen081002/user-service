package domain

import (
	"context"
	"erp/internal/api/request"
	"erp/internal/infrastructure"
	"erp/internal/models"
)

type UserRepository interface {
	GetByID(db *infrastructure.Database, ctx context.Context, id string) (res *models.User, err error)
	IsExistEmail(db *infrastructure.Database, ctx context.Context, email string) (res *models.User, err error)
	Create(db *infrastructure.Database, ctx context.Context, user *models.User) (res *models.User, err error)
	ListUsers(db infrastructure.Database, o request.PageOptions, ctx context.Context) (res []*models.User, total *int64, err error)
	GetByEmail(db *infrastructure.Database, ctx context.Context, email string) (res *models.User, err error)
	UpdateLastLogin(db *infrastructure.Database, ctx context.Context, id string) error
}

type UserService interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(o request.PageOptions, ctx context.Context) (res []*models.User, total *int64, err error)
}

type ListUsersRequest struct {
	request.PageOptions
}
