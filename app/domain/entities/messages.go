package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

//Meesage cotains data body
type Message struct {
	ID      primitive.ObjectID `bson:"_id,omitemtpy" json:"id,omitempty"`
	Topic   string             `bson:"topic" json:"topic" validate:"required,min=3,max=5"`
	Payload []byte             `bson:"payload" json:"payload" validate:"required,max=1000"`
}
