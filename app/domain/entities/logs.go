package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Log saves messages failed
type Logs struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Desc       string             `bson:"description" json:"description"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	Stage      string             `bson:"stage" json:"stage"`
	Subscriber string             `bson:"subscriber" json:"subscriber"`
}
