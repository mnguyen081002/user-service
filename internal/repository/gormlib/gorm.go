package gormlib

import (
	"erp/internal/api/request"
	"erp/internal/infrastructure"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type GormRepository struct {
	infrastructure.Database
}

func NewGormTransaction(db infrastructure.Database) infrastructure.DatabaseTransaction {
	return GormRepository{db}
}

func (g GormRepository) WithTransaction(txFunc func(tx *infrastructure.Database) error) (err error) {
	err = g.RDBMS.Transaction(func(tx *gorm.DB) error {
		g.RDBMS = tx
		err := txFunc(&g.Database)
		return err
	})
	return err
}

type GormQueryPaginationBuilder[E any] struct {
	tx *gorm.DB
}

func GormQueryPagination[E any](tx *gorm.DB, o request.PageOptions, data *[]*E) *GormQueryPaginationBuilder[E] {
	q := &GormQueryPaginationBuilder[E]{
		tx: tx,
	}
	if o.Page == 0 {
		o.Page = 1
	}
	if o.Limit == 0 {
		o.Limit = 10
	}
	offset := (o.Page - 1) * o.Limit

	q.tx = q.tx.Debug().Offset(int(offset)).Limit(int(o.Limit)).Find(&data)
	return q
}

func (q *GormQueryPaginationBuilder[E]) Count(total *int64) *GormQueryPaginationBuilder[E] {
	if total == nil {
		total = new(int64)
	}
	q.tx.Count(total)
	return q
}

func (q *GormQueryPaginationBuilder[E]) Error() error {
	return errors.WithStack(q.tx.Error)
}
