package models

import (
	"erp/internal/constants"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type User struct {
	BaseModel   `bson:",inline"`
	FirstName   string         `json:"first_name" gorm:"column:first_name;type:varchar(50);not null" bson:"first_name"`
	LastName    string         `json:"last_name" gorm:"column:last_name;type:varchar(50);not null" bson:"last_name"`
	Email       string         `json:"email" gorm:"column:email;type:varchar(100);not null;index:idx_email,unique" bson:"email"`
	Password    string         `json:"password" gorm:"column:password;type:varchar(255);not null" bson:"password"`
	Role        constants.Role `json:"role" gorm:"column:role;type:varchar(50);not null" bson:"role"`
	LastLoginAt time.Time      `json:"last_login_at" gorm:"column:last_login_at;type:timestamp" bson:"last_login_at"`
}

func (u *User) MarshalBSON() ([]byte, error) {
	if u.CreatedAt.IsZero() {
		u.ID = uuid.NewV4()
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()

	type my User
	return bson.Marshal((*my)(u))
}
