package mongo

import (
	"context"
	"erp/internal/api/request"
	"erp/internal/infrastructure"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	infrastructure.Database
}

func NewMongoTransaction(db infrastructure.Database) infrastructure.DatabaseTransaction {
	return MongoRepository{db}
}

func (g MongoRepository) WithTransaction(txFunc func(tx *infrastructure.Database) error) (err error) {
	err = g.Mongo.Client().UseSession(context.Background(), func(sessionContext mongo.SessionContext) error {
		err := sessionContext.StartTransaction()
		if err != nil {
			return errors.WithStack(err)
		}
		g.Context = sessionContext
		err = txFunc(&g.Database)
		if err != nil {
			err1 := sessionContext.AbortTransaction(sessionContext)
			if err1 != nil {
				return errors.WithStack(err1)
			}
			return errors.WithStack(err)
		}
		err = sessionContext.CommitTransaction(sessionContext)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

type MongoQueryPaginationBuilder[E any] struct {
	collection *mongo.Collection
	filter     interface{}
	err        error
}

func MongoQueryPagination[E any](collection *mongo.Collection, filter interface{}, o request.PageOptions, data *[]*E) *MongoQueryPaginationBuilder[E] {
	q := &MongoQueryPaginationBuilder[E]{
		collection: collection,
		filter:     filter,
	}
	if o.Page == 0 {
		o.Page = 1
	}
	if o.Limit == 0 {
		o.Limit = 10
	}
	offset := (o.Page - 1) * o.Limit
	c, err := q.collection.Find(context.Background(), filter, &options.FindOptions{
		Limit: &o.Limit,
		Skip:  &offset,
	})

	if err != nil {
		q.err = errors.WithStack(err)
		return q
	}

	err = c.All(context.Background(), data)

	if err != nil {
		q.err = errors.WithStack(err)
	}

	return q
}

func (q *MongoQueryPaginationBuilder[E]) Count(total *int64) *MongoQueryPaginationBuilder[E] {
	if q.err != nil {
		return q
	}
	t, err := q.collection.CountDocuments(context.Background(), q.filter, &options.CountOptions{})
	if err != nil {
		q.err = errors.WithStack(err)
	}

	*total = t
	return q
}

func (q *MongoQueryPaginationBuilder[E]) Error() error {
	return q.err
}
