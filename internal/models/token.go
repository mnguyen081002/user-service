package models

import (
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Token struct {
	BaseModel `bson:",inline"`
	UserID    string `json:"user_id" gorm:"column:user_id;type:uuid;not null;index:idx_user_id,unique" bson:"user_id"`
	ExpiredAt int64  `json:"expired_at" gorm:"column:expired_at;type:bigint;not null" bson:"expired_at"`
}

func (u *Token) MarshalBSON() ([]byte, error) {
	if u.CreatedAt.IsZero() {
		u.ID = uuid.NewV4()
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()

	type my Token
	return bson.Marshal((*my)(u))
}
