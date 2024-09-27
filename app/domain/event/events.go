package event

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/topic"
	"dominus/app/interfaces/clients"

	"go.mongodb.org/mongo-driver/bson"
)

type events struct {
	t topic.TopicInt
	r clients.ClientInt
}

func (e events) InitialLoad() {
	var (
		repo = e.r.NewMongoClient()
		topics  []entities.TopicDB
	)

	if err := repo.FindObjects("topic", &topics, bson.D{}); err != nil {
		return
	}

	for _, it := range topics {
		e.t.CreateTopic(it.Topic, it.Subscribers)
	}
}

func (e events) Resend() {

}





type EventsInt interface {
	//Initial load to feed topic,looks for information on topic collection
	InitialLoad()
	//Resend patter sends message failed
	Resend()
}



func NewEvent(rep clients.ClientInt, t topic.TopicInt) EventsInt {
	return &events{t: t, r: rep}
}