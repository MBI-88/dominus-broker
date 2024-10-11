package event

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/topic"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type events struct {
	t      topic.TopicInt
	repo   repoInt
	rest   restClientInt
	status bool
}

func (e *events) InitialLoad() {
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

func NewEvent(r repoInt, rs restClientInt, t topic.TopicInt) EventsInt {
	return &events{
		t:      t,
		repo:   r,
		rest:   rs,
		status: false,
	}
}


type repoInt interface {
	//Find objects in the database
	//
	//Parameters
	//
	//-> collection: collection name to use
	//
	//-> objects: array object to fill
	//
	//-> filter: the filter to find objects
	//
	//-> op: contains options to use in the query
	FindObjects(collection string, objects any, filter bson.D, op ...*options.FindOptions) error
}

type restClientInt interface {
	//Rest client
	//
	//Parameters
	//
	//-> ctx: context
	// 
	//-> sub: subcriber
	//
	//-> payload: message
	DoJsonRequest(sub string, payload []byte) error
}
