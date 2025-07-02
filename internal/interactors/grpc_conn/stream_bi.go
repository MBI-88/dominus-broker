package grpcconn

import (
	"context"
	"dominus-project/internal/domain/repos"
	"fmt"
)

func (c *grpcService) StreamBiConn(stream repos.StreamBiInt) error {
	closed := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := stream.Recv()
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("Subscribers not found")
	}
	total := len(subscribers)
	done := make(chan struct{}, total)
	streamProv := make(chan []byte)
	streamSub := make(chan []byte, total+int(total*2/3))
	errMsg := make(chan error, total)

	if err != nil {
		return err
	}
	go func(sig <-chan error) {
		for {
			select {
			case err, ok := <-sig:
				if ok {
					go c.lg.WriteLog("InsertObject", err.Error())
				} else {
					return
				}
			}
		}
	}(errMsg)

	go c.client.BidirectionalStream(subscribers, streamProv, streamSub, errMsg, closed, ctx, done)
	streamProv <- req.GetPayload()

	//Receives from provider
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				go c.lg.WriteLog("StreamBiConn", err.Error())
				close(streamProv)
				<-closed
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
					go c.lg.WriteLog("StreamBiConn", err.Error())
					cancel()
				}
			}
		case <-done:
			total--
			if total == 0 {
				close(streamSub)
				close(errMsg)
				close(done)
				return fmt.Errorf("Connection closed")
			}
		}
	}
}
