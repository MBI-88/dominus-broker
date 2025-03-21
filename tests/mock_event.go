package tests

import (
	"dominus-project/app/domain/repos"
	"fmt"
	"log"
)



type eventMock struct {
	statusFlag bool
}

func (*eventMock) WriteLog(op, dsc string) {
	log.Println(op, " ", dsc)
}

func (*eventMock) Printf(format string, args ...any) {
	log.Println(format, " ", args)
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