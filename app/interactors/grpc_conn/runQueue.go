package grpcconn

import "time"

func (c *grpcService) RunQueue() {
	for {
		select {
		case <-time.Tick(2 * time.Second):
			c.checkQueue()
		}
	}
}
