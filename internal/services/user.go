package service

import (
	"context"
	config "erp/config"
	"erp/internal/api/request"
	"erp/internal/domain"
	"erp/internal/infrastructure"
	"erp/internal/models"
	"erp/internal/repository"
)

type (
	userService struct {
		ufw    *repository.UnitOfWork
		db     infrastructure.Database
		config *config.Config
	}
)

func NewUserService(ufw *repository.UnitOfWork, db infrastructure.Database, config *config.Config) domain.UserService {
	return &userService{
		ufw:    ufw,
		db:     db,
		config: config,
	}
}
func (u *userService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if _, err := u.ufw.UserRepository.Create(&u.db, ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userService) GetByID(ctx context.Context, id string) (user *models.User, err error) {
	user, err = u.ufw.UserRepository.GetByID(&u.db, ctx, id)
	return
}

func (u *userService) ListUsers(o request.PageOptions, ctx context.Context) (res []*models.User, total *int64, err error) {
	res, total, err = u.ufw.UserRepository.ListUsers(u.db, o, ctx)
	return
}

func (u *userService) GetByEmail(ctx context.Context, email string) (user *models.User, err error) {
	user, err = u.ufw.UserRepository.GetByEmail(&u.db, ctx, email)
	return
}
