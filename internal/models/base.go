package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type BaseModel struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()" bson:"_id"`
	UpdaterID *uuid.UUID `json:"updater_id" gorm:"column:updater_id;type:uuid;" bson:"updater_id"`
	CreatedAt time.Time  `gorm:"column:created_at;type:timestamp;default:now();not null" json:"created_at" bson:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;type:timestamp;default:now();not null" json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:timestamp" json:"deleted_at" bson:"deleted_at"`
}
