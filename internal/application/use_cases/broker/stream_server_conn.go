package broker

import (
	"context"

	"dominus-broker/internal/application/dtos"
	"fmt"
)

func (b *broker) StreamServerConn(req dtos.BrokerRequestDto, st dtos.BrokerServerDto) error {
	ctx, cancel := context.WithCancel(st.Context())
	defer cancel()
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("subscribers not found")
	}
	total := len(subscribers)
	stream := make(chan []byte, total+int(total*2/3))
	closed := make(chan struct{})
	initialRequest := req.GetPayload()

	go b.client.ServerStream(subscribers, initialRequest, stream, ctx, closed)

	for {
		select {
		case body, ok := <-stream:
			if ok {
				if err := st.Send(body); err != nil {
					cancel()
				}
			}
		case <-closed:
			close(closed)
			close(stream)
			return fmt.Errorf("connection closed")
		}
	}
}
