package event

import (
	"dominus/app/domain/topic"
	"dominus/app/interfaces/clients"
)

type events struct {
	t topic.TopicInt
	r clients.ClientInt
}

func (e events) InitialLoad() {

}

func (e events) Resend() {

}





type EventsInt interface {
	//Initial load to feed topic
	InitialLoad()
	//Resend patter for message send resiliency
	Resend()
}



func NewEvent(rep clients.ClientInt, t topic.TopicInt) EventsInt {
	return &events{t: t, r: rep}
}