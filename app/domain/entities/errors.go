package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Log saves messages failed
type Logs struct {
	ID        primitive.ObjectID `bson:"_id,omitemtpy" json:"id,omitempty"`
	Desc      string             `bson:"description" json:"description"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	Status    uint32             `bson:"status" json:"status"`
	Sub       string             `bson:"sub" josn:"sub"`
}
