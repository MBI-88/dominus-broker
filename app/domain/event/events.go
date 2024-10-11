package event

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/topic"
	"dominus/app/interfaces/database"
	"dominus/app/interfaces/rest/output"

	"go.mongodb.org/mongo-driver/bson"
)

type events struct {
	t      topic.TopicInt
	repo   database.RepositoryInt
	rest   output.RestClientInt
	status bool
}

func (e events) InitialLoad() {
	var topics []entities.Topic

	if err := e.repo.FindObjects("topic", &topics, bson.D{}); err != nil {
		return
	}

	for _, it := range topics {
		e.t.CreateTopic(it.Topic, it.Subscribers)
	}
}

func (e events) Sentinel(status <-chan bool) {
	for {
		select {
		case r, ok := <-status:
			if ok && r && !e.status {
				go e.resend()
			}
			if ok == false {
				break
			}
		}
	}
}

// if there are any pending message , resend must stop
func (e *events) resend() {
	e.status = true

}

type EventsInt interface {
	//Initial load to feed topic,looks for information on topic collection
	InitialLoad()
	//Check client fails
	//
	//Parameters
	//
	//-> status: channel
	//
	//-> buffer: total gorutine alive
	Sentinel(status <-chan bool)
}

func NewEvent(r database.RepositoryInt, rs output.RestClientInt, t topic.TopicInt) EventsInt {
	return &events{
		t:      t,
		repo:   r,
		rest:   rs,
		status: false,
	}
}
