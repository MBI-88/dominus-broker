package sqs

import (
	"context"
	"fmt"
	"os"
	"time"
)

func (q *sqs) ReactivateMessage(ch <-chan os.Signal) {
	go func() {
		for {
			select {
			case <-ch:
				return
			case <-time.Tick(time.Minute):
				ctx := context.Background()
				for i := 0; i < q.memory.Len(); i++ {
					key := q.memory.Iter(i)

					payload, err := q.client.GetMessage(ctx, key)
					if err != nil {
						continue
					}

					if payload.GetHidden() && time.Since(payload.GetCreatedAt()).Minutes() > 1 {
						payload.SetHidden(false)
						payload.SetCreateAt(time.Now())

						if err := q.client.SendMessage(ctx, payload); err != nil {
							fmt.Printf("%s\n", err)
						}
					}
				}
			}
		}
	}()
}
