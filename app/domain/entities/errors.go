package entities

import (
	"time"
)

//Log saves messages failed
type Logs struct {
	Log       Message            `bson:"inline"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	Status    string             `bson:"status" json:"status"`          // sent | pending
} 
