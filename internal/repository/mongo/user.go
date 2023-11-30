package mongo

import (
	"context"
	"erp/internal/api/request"
	"erp/internal/api_errors"
	"erp/internal/domain"
	"erp/internal/infrastructure"
	"erp/internal/models"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

type userRepositoryImpl struct {
	logger         *zap.Logger
	collectionName string
}

func NewUserRepository() domain.UserRepository {
	return userRepositoryImpl{
		collectionName: "users",
	}
}

func (u userRepositoryImpl) GetByID(db *infrastructure.Database, ctx context.Context, id string) (res *models.User, err error) {
	err = db.Mongo.Collection(u.collectionName).FindOne(ctx, id).Decode(&res)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func (u userRepositoryImpl) IsExistEmail(db *infrastructure.Database, ctx context.Context, email string) (res *models.User, err error) {
	err = db.Mongo.Collection(u.collectionName).FindOne(ctx, email).Decode(&res)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func (u userRepositoryImpl) Create(db *infrastructure.Database, ctx context.Context, user *models.User) (res *models.User, err error) {
	_, err = db.Mongo.Collection(u.collectionName).InsertOne(ctx, user)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return user, nil
}

func (u userRepositoryImpl) ListUsers(db infrastructure.Database, o request.PageOptions, ctx context.Context) ([]*models.User, *int64, error) {
	var res []*models.User
	total := new(int64)

	err := MongoQueryPagination(db.Mongo.Collection(u.collectionName), bson.M{}, o, &res).Count(total).Error()
	if err != nil {
		return nil, nil, err
	}

	return res, total, nil
}

func (u userRepositoryImpl) GetByEmail(db *infrastructure.Database, ctx context.Context, email string) (res *models.User, err error) {
	err = db.Mongo.Collection(u.collectionName).FindOne(db.Context, bson.M{
		"email": email,
	}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(api_errors.ErrEmailNotFound)
		}
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func (u userRepositoryImpl) UpdateLastLogin(db *infrastructure.Database, ctx context.Context, id string) error {
	_, err := db.Mongo.Collection(u.collectionName).UpdateOne(db.Context, bson.M{
		"_id": uuid.FromStringOrNil(id),
	}, bson.M{
		"$set": bson.M{
			"last_login_at": time.Now(),
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
