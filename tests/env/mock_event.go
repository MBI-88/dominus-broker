package env

import (
	"dominus-project/internal/domain/repos"
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
		return make([]string, 10), nil
	}
	return nil, fmt.Errorf("[-] Error response")
}

func NewEventMock(status bool) repos.ILogs {
	return &eventMock{
		statusFlag: status,
	}
}
