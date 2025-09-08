package grpcconn

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"fmt"
)

func (c *grpcService) StreamServerConn(req adapters.GrpcDto, st adapters.StreamServer) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("subscribers not found")
	}
	total := len(subscribers)
	stream := make(chan []byte, total+int(total*2/3))
	done := make(chan struct{}, total)
	initialRequest := req.GetPayload()

	go c.client.ServerStream(subscribers, initialRequest, stream, ctx, done)

	for {
		select {
		case body, ok := <-stream:
			if ok {
				if err := st.Send(body); err != nil {
					go c.lg.WriteLog("StreamServerConn", err.Error())
					cancel()
				}
			}
		case <-done:
			total--
			if total == 0 {
				close(done)
				close(stream)
				return fmt.Errorf("connection closed")
			}
		}
	}
}
