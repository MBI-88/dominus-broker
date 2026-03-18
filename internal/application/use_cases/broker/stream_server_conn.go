package broker

import (
	"context"

	"dominus-project/internal/application/dtos"
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
	done := make(chan struct{}, total)
	initialRequest := req.GetPayload()

	go b.client.ServerStream(subscribers, initialRequest, stream, ctx, done)

	for {
		select {
		case body, ok := <-stream:
			if ok {
				if err := st.Send(body); err != nil {
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
