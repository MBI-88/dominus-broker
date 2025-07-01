package grpcconn

import (
	"fmt"
	"time"
)

func (c *grpcService) RunQueue(close <-chan struct{}) error {
	for {
		select {
		case <-time.Tick(1 * time.Second):
			c.checkQueue()
		case <-close:
			return fmt.Errorf("Queue closed")
		}
	}
}
