package env

import (
	"context"
	"dominus-project/internal/domain/repositories"
	"log"
)

type eventMock struct {
	statusFlag bool
}

func NewEventMock(status bool) repositories.Logs {
	return &eventMock{
		statusFlag: status,
	}
}

func (*eventMock) CheckID(ctx context.Context) context.Context {
	return  ctx
}

func (*eventMock) WriteLog(ctx context.Context, level, op, dsc string) {
	log.Println(op, " ", dsc)
}

