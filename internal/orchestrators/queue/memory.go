package queue

import (
	"context"
	"dominus-project/internal/domain/enum"
)

func (q *queue) CheckMemory() {
	ctx := q.lg.CheckID(context.Background())
	if !q.memory.CheckMemory() {
		q.lg.WriteLog(ctx, enum.ERROR, "CheckMemory", "memory empty")
		if err := q.client.GetKeys(ctx, q.memory); err != nil {
			q.lg.WriteLog(ctx, enum.ERROR, "CheckMemory", err.Error())
		}
		q.lg.WriteLog(ctx, enum.INFO, "CheckID", "memory ready")
		return
	}

	q.lg.WriteLog(ctx, enum.INFO, "CheckMemory" ,"memory loaded")
}