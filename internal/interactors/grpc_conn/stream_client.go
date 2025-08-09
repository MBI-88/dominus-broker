package grpcconn

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"fmt"
)

func (c *grpcService) StreamClientConn(st adapters.StreamClient) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := make(chan []byte)
	req, err := st.Recv()
	if err != nil {
		return err
	}
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("Subscribers not found")
	}

	go c.client.ClientStream(subscribers, stream, ctx)

	stream <- req.GetPayload()
	for {
		req, err := st.Recv()
		if err != nil {
			c.lg.WriteLog("StreamClientConn", err.Error())
			close(stream)
			cancel()
			return err
		}
		stream <- req.GetPayload()
	}
}
