package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Topic to serializer in the database
type Topic struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Topic       string             `bson:"topic,omitempty" json:"topic,omitempty" validate:"omitempty,min=3,max=10,alphanum"`
	Subscribers []string           `bson:"subscribers,omitempty" json:"subscribers,omitempty" validate:"omitempty,dive,url,max=10"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
