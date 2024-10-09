package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"dominus/app/domain/topic"
	"dominus/app/interfaces/clients"

	jsoniter "github.com/json-iterator/go"
)

type interactor struct {
	client clients.ClientInt
	topic  topic.TopicInt
	event  event.EventsInt
	log    event.LogsInt
}

func (i interactor) NewConnection() ConnectionInt {
	return &connection{
		m:  new(entities.Message),
		p:  jsoniter.ConfigCompatibleWithStandardLibrary,
		t:  i.topic,
		ev: i.event,
		r:  rules.NewRule(),
		lg: i.log,
	}
}

func (i interactor) NewManager() ManagerInt {
	return &manager{
		p:          jsoniter.ConfigCompatibleWithStandardLibrary,
		r:          rules.NewRule(),
		t:          i.topic,
		repo:       i.client.MongoClient(),
		collection: "topic",
		lg:         i.log,
	}
}

type InteractorInt interface {
	NewConnection() ConnectionInt
	NewManager() ManagerInt
}

// Create a new interactor instance
func NewInteractor(ct clients.ClientInt, t topic.TopicInt, e event.EventsInt, lg event.LogsInt) InteractorInt {
	return &interactor{
		topic:  t,
		event:  e,
		client: ct,
		log:    lg,
	}
}
