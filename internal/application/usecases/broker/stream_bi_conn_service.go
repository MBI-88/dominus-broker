package broker

import (
	"context"

	"fmt"
)

func (b *broker) StreamBiConn(stream BrokerBidirectionalDto) error {
	closed := make(chan struct{})
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("broker.StreamBiConn %s", err)
	}
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("broker.StreamBiConn subscribers not found")
	}
	total := len(subscribers)
	streamProv := make(chan []byte)
	streamSub := make(chan []byte, total+int(total*2/3))

	go b.client.BidirectionalStream(subscribers, streamProv, streamSub, closed, ctx)
	streamProv <- req.GetPayload()

	//Receives from provider
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				close(streamProv)
				return
			}
			streamProv <- req.GetPayload()
		}
	}()

	//Receives from subscribers
	for {
		select {
		case payload, ok := <-streamSub:
			if ok {
				if err := stream.Send(payload); err != nil {
					cancel()
				}
			}
		case <-closed:
			close(streamSub)
			close(closed)
			return fmt.Errorf("broker.StreamBiConn connection closed")
		}
	}
}
