package entities

import (
	"time"
)

// Log saves messages failed
type Logs struct {
	ID         string    `bson:"_id,omitemtpy" json:"id,omitempty"`
	Desc       string    `bson:"description" json:"description"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
	Status     uint32    `bson:"status" json:"status"`
	Stage      string    `bson:"stage" json:"stage"`
	Subscriber string    `bson:"subscriber" json:"subscriber"`
}
