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
	topic topic.TopicInt
	event event.EventsInt
}

func (i interactor) NewConnection() ConnectionInt {
	return  &connection{
		m: new(entities.Message),
		p: jsoniter.ConfigCompatibleWithStandardLibrary,
		t: i.topic,
		ev: i.event,
	}
}

func (i interactor) NewManager() ManagerInt {
	return &manager{
		p: jsoniter.ConfigCompatibleWithStandardLibrary,
		r: rules.NewRule(),
		t: i.topic,
		repo: i.client.NewMongoClient(),
	}
}


type InteractorInt interface {
	NewConnection() ConnectionInt 
	NewManager() ManagerInt
}



// Create a new interactor instance
func NewInteractor(ct clients.ClientInt, t topic.TopicInt, e event.EventsInt) InteractorInt {
	return &interactor{
		topic: t,
		event: e,
		client: ct,
	}
}
