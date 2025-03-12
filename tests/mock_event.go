package tests

import (
	"dominus-project/app/domain/repos"
	"fmt"
)



type eventMock struct {
	statusFlag bool
}

func (*eventMock) WriteLog(op, dsc string) {

}

func (*eventMock) Printf(format string, args ...any) {
	
}

func (ev *eventMock) GetLogs() ([]string, error) {
    if ev.statusFlag {
		return nil, nil
	}
	return nil, fmt.Errorf("[-] Error response")
}



func NewEventMock(status bool) repos.LogsInt {
	return &eventMock{
		statusFlag: status,
	}
}