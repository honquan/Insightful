package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type Mongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt int64              `bson:"c_at,omitempty"`
	UpdatedAt int64              `bson:"u_at,omitempty"`
	//UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}

type MySQL struct {
	gorm.Model
}

type MongoModel interface {
	CollectionName(ws string) string
}

type SQLModel interface {
	TableName() string
}
