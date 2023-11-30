package mongo

import (
	"context"
	"erp/internal/domain"
	"erp/internal/infrastructure"
	"erp/internal/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// for mongodb driver

type tokenRepositoryImpl struct {
	collectionName string
}

func NewTokenRepository() domain.TokenRepository {
	return tokenRepositoryImpl{
		collectionName: "tokens",
	}
}

func (t tokenRepositoryImpl) Upsert(db *infrastructure.Database, ctx context.Context, token *models.Token) (res *models.Token, err error) {
	_, err = db.Mongo.Collection(t.collectionName).UpdateOne(db.Context, bson.M{
		"user_id": token.UserID,
	}, bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}, options.Update().SetUpsert(true))

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return token, nil
}
