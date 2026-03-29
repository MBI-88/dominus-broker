package broker

import (
	"context"

	"dominus-project/internal/application/dtos"
	"fmt"
)

func (b *broker) StreamBiConn(stream dtos.BrokerBidirectionalDto) error {
	closed := make(chan struct{})
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("subscribers not found")
	}
	total := len(subscribers)
	done := make(chan struct{}, total)
	streamProv := make(chan []byte)
	streamSub := make(chan []byte, total+int(total*2/3))
	errMsg := make(chan error, total)

	go func(sig <-chan error) {
		for er := range sig {
			fmt.Println(er)
		}
	}(errMsg)

	go b.client.BidirectionalStream(subscribers, streamProv, streamSub, errMsg, closed, ctx, done)
	streamProv <- req.GetPayload()

	//Receives from provider
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				close(streamProv)
				<-closed
				close(closed)
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
		case <-done:
			total--
			if total == 0 {
				close(streamSub)
				close(errMsg)
				close(done)
				return fmt.Errorf("connection closed")
			}
		}
	}
}
