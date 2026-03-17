package queue

import (
	"context"
	"fmt"
)

func (q *queue) CheckMemory() {
	if !q.memory.CheckMemory() {
		if err := q.client.GetKeys(context.Background(), q.memory); err != nil {
			fmt.Println(err)
		}
		return
	}

	fmt.Println("Memory load completed")
}