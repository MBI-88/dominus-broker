package tests

import (
	"dominus/app/domain/repos"
)



type eventMock struct {
}

func (*eventMock) WriteLog(op, dsc string) {

}

func (*eventMock) Printf(format string, args ...any) {
	
}

func (*eventMock) GetLogs() ([]string, error) {

	return nil, nil
}



func NewEventMock() repos.LogsInt {
	return new(eventMock)
}