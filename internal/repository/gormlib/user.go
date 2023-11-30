package gormlib

import (
	"context"
	"erp/internal/api/request"
	"erp/internal/api_errors"
	"erp/internal/domain"
	"erp/internal/infrastructure"
	"erp/internal/models"
	"erp/internal/utils"
	errors2 "errors"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type userRepositoryImpl struct {
	logger *zap.Logger
}

func NewUserRepository() domain.UserRepository {
	return userRepositoryImpl{}
}

func (u userRepositoryImpl) GetByEmail(db *infrastructure.Database, ctx context.Context, email string) (res *models.User, err error) {
	err = db.RDBMS.WithContext(ctx).Where("email = ?", email).First(&res).Error
	if err != nil {
		if utils.ErrNoRows(err) {
			return res, errors.New(api_errors.ErrUserNotFound)
		}
		return nil, err
	}
	return
}

func (u userRepositoryImpl) Create(db *infrastructure.Database, ctx context.Context, user *models.User) (res *models.User, err error) {
	err = db.RDBMS.Create(&user).Error
	if errors2.Is(err, gorm.ErrDuplicatedKey) {
		return nil, errors.New(api_errors.ErrEmailAlreadyExists)
	}
	return user, err
}

func (u userRepositoryImpl) GetByID(db *infrastructure.Database, ctx context.Context, id string) (res *models.User, err error) {
	err = db.RDBMS.WithContext(ctx).Where("id = ?", id).First(&res).Error
	if err != nil {
		if utils.ErrNoRows(err) {
			return res, errors.New(api_errors.ErrUserNotFound)
		}
		return nil, err
	}
	return
}

func (u userRepositoryImpl) IsExistEmail(db *infrastructure.Database, ctx context.Context, email string) (res *models.User, err error) {
	err = db.RDBMS.WithContext(ctx).Where("email = ?", email).First(&res).Error
	if err != nil {
		if utils.ErrNoRows(err) {
			return res, errors.New(api_errors.ErrUserNotFound)
		}
		return nil, err
	}
	return
}

func (u userRepositoryImpl) ListUsers(db infrastructure.Database, o request.PageOptions, ctx context.Context) ([]*models.User, *int64, error) {
	var res = make([]*models.User, 0)
	var total = new(int64)

	q := db.RDBMS.WithContext(ctx).Model(&models.User{})

	err := GormQueryPagination(q, o, &res).Count(total).Error()
	if err != nil {
		return nil, nil, err
	}
	return res, total, nil
}

func (u userRepositoryImpl) UpdateLastLogin(db *infrastructure.Database, ctx context.Context, id string) error {
	return db.RDBMS.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("last_login_at", time.Now()).Error
}
