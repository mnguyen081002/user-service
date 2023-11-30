package utils

import (
	"erp/internal/infrastructure"
	"errors"
	"gorm.io/gorm"
)

func ErrNoRows(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func MustHaveDb(db *infrastructure.Database) {
	if db == nil {
		panic("Database engine is null")
	}
}
