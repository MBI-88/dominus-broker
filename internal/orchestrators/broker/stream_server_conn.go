package broker

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/enum"
	"fmt"
)

func (b *broker) StreamServerConn(req adapters.BrokerRequestDto, st adapters.BrokerServerDto) error {
	ctx, cancel := context.WithCancel(st.Context())
	defer cancel()
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		go b.lg.WriteLog(ctx, enum.ERROR, "StreamServerConn", "Subscribers not found")
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
					go b.lg.WriteLog(ctx, enum.ERROR ,"StreamServerConn", err.Error())
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
