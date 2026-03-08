package broker

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/enum"
	"fmt"
)

func (b *broker) StreamBiConn(stream adapters.BrokerBidirectionalDto) error {
	closed := make(chan struct{})
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	req, err := stream.Recv()
	if err != nil {
		go b.lg.WriteLog(ctx, enum.ERROR, "StreamBiConn", err.Error())
		return err
	}
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		go b.lg.WriteLog(ctx, enum.ERROR, "StreamBiConn", enum.SUBCRIBER_NOT_FOUND)
		return fmt.Errorf("subscribers not found")
	}
	total := len(subscribers)
	done := make(chan struct{}, total)
	streamProv := make(chan []byte)
	streamSub := make(chan []byte, total+int(total*2/3))
	errMsg := make(chan error, total)

	go func(sig <-chan error) {
		for er := range sig {
			go b.lg.WriteLog(ctx, enum.ERROR, "InsertObject", er.Error())
		}
	}(errMsg)

	go b.client.BidirectionalStream(subscribers, streamProv, streamSub, errMsg, closed, ctx, done)
	streamProv <- req.GetPayload()

	//Receives from provider
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				go b.lg.WriteLog(ctx, enum.ERROR, "StreamBiConn", err.Error())
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
					go b.lg.WriteLog(ctx, enum.ERROR, "StreamBiConn", err.Error())
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
