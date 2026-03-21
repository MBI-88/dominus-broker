package sqs

import (
	"context"
	"fmt"
)

func (q *sqs) CheckMemory() {
	if !q.memory.CheckMemory() {
		if err := q.client.GetKeys(context.Background(), q.memory); err != nil {
			fmt.Println(err)
		}
		return
	}

	fmt.Println("Memory load completed")
}