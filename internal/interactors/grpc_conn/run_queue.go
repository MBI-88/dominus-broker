package grpcconn

import (
	"fmt"
	"time"
)

func (c *grpcService) RunQueue(close <-chan struct{}) error {
	for {
		select {
		case <-time.Tick(600 * time.Millisecond):
			go c.checkQueue()
		case <-close:
			return fmt.Errorf("queue closed")
		}
	}
}
