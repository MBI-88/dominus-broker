package interactors

import (
	"context"
	"dominus-project/app/domain/repos"
	"fmt"
)

func (c *grpcService) StreamClientConn(st repos.StreamClientInt) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := make(chan []byte, 0)
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

func (c *grpcService) StreamServerConn(req repos.GrpRequestMessageInt, st repos.StreamServerInt) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("Subscribers not found")
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
				return fmt.Errorf("Connection closed")
			}
		}
	}
}

func (c *grpcService) StreamBiConn(stream repos.StreamBiInt) error {
	closed := make(chan struct{}, 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req, err := stream.Recv()
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("Subscribers not found")
	}
	total := len(subscribers)
	done := make(chan struct{}, total)
	streamProv := make(chan []byte, 0)
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