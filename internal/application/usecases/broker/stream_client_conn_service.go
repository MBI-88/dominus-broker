package broker

import (
	"context"

	"fmt"
)

func (b *broker) StreamClientConn(st BrokerClientDto) error {
	ctx, cancel := context.WithCancel(st.Context())
	defer cancel()
	stream := make(chan []byte)
	req, err := st.Recv()
	if err != nil {
		return err
	}
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("broker.StreamClientConn subscribers not found")
	}

	go b.client.ClientStream(subscribers, stream, ctx)

	stream <- req.GetPayload()
	for {
		req, err := st.Recv()
		if err != nil {
			close(stream)
			cancel()
			return err
		}
		stream <- req.GetPayload()
	}
}
