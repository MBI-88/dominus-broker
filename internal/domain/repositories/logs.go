package repositories

import (
	"context"
)

type Logs interface {
	WriteLog(ctx context.Context, level, op, dsc string)
	CheckID(ctx context.Context) context.Context
}
