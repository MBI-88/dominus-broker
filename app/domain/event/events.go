package event



type events struct {
}

func (events) InitialLoad() {

}

func (events) Resend() {

}

type EventsInt interface {
	InitialLoad()
	Resend()
}



func NewEvent() EventsInt {
	return &events{}
}