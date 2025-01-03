package tests

import "dominus/app/domain/event"



type eventMock struct {
}

func (*eventMock) WriteLog(op, dsc string) {

}

func (*eventMock) Printf(format string, args ...any) {
	
}



func NewEventMock() event.LogsInt {
	return new(eventMock)
}