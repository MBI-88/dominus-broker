package tests

import "dominus/app/interactors"



type eventMock struct {
}

func (*eventMock) WriteLog(op, dsc string) {

}

func (*eventMock) Printf(format string, args ...any) {
	
}



func NewEventMock() interactors.LogsInt {
	return new(eventMock)
}