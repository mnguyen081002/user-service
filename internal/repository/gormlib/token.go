package gormlib

import (
	"context"
	"erp/internal/domain"
	"erp/internal/infrastructure"
	"erp/internal/models"
	"gorm.io/gorm/clause"
)

type tokenRepositoryImpl struct {
}

func NewTokenRepository() domain.TokenRepository {
	return tokenRepositoryImpl{}
}

func (t tokenRepositoryImpl) Upsert(db *infrastructure.Database, ctx context.Context, token *models.Token) (res *models.Token, err error) {
	if err = db.RDBMS.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		UpdateAll: true,
	}).Create(&token).Error; err != nil {
		return nil, err
	}
	return token, nil
}
