package interactors

import (
	"dominus/app/domain/repos"
	"dominus/app/domain/rules"
	"fmt"
	"time"
)

type grpcService struct {
	rls        rules.RulesInt
	client     repos.GrpClientInt
	lg         repos.LogsInt
	collection string
}

func (c *grpcService) SimpleConn(ms repos.GrpRequestMessageInt) error {
	subs := ms.GetSubscribers()
	if len(subs) == 0 {
		return fmt.Errorf("Subscribers not found")
	}
	body := ms.GetPayload()
	for _, sub := range subs {
		go func(url string, body []byte) {
			_, err := c.client.Simple(url, body)
			if err != nil {
				c.lg.WriteLog("SimpleConn", err.Error())
			}
		}(sub, body)
	}
	return nil
}

func (c *grpcService) StreamClientConn(st repos.StreamClientInt) error {
	stream := make(chan []byte, 0)
	closed := make(chan struct{}, 0)
	req, err := st.Recv()
	if err != nil {
		return err
	}
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("Subscribers not found")
	}

	errMsg := make(chan error, len(subscribers))
	go c.client.ClientStream(subscribers, stream, errMsg, closed)
	go func(sig <-chan error) {
		for {
			select {
			case err, ok := <-sig:
				if ok {
					c.lg.WriteLog("StreamClientConn", err.Error())
				} else {
					return
				}
			}
		}
	}(errMsg)

	stream <- req.GetPayload()
	for {
		req, err := st.Recv()
		if err != nil {
			c.lg.WriteLog("StreamClientConn", err.Error())
			close(stream)
			<-closed
			close(errMsg)
			return err
		}
		stream <- req.GetPayload()
	}
}

func (c *grpcService) StreamServerConn(req repos.GrpRequestMessageInt, st repos.StreamServerInt) error {
	closed := make(chan struct{}, 0)
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("Subscribers not found")
	}
	total := len(subscribers)
	stream := make(chan []byte, total+int(total*2/3))
	errMsg := make(chan error, total)
	done := make(chan struct{}, total)
	initialRequest := req.GetPayload()

	go c.client.ServerStream(subscribers, initialRequest, stream, errMsg, closed, done)

	go func(sig <-chan error) {
		for {
			select {
			case err, ok := <-sig:
				if ok {
					go c.lg.WriteLog("StreamServerConn", err.Error())
				} else {
					return
				}
			}
		}
	}(errMsg)

	for {
		select {
		case body, ok := <-stream:
			if ok {
				if err := st.Send(body); err != nil {
					go c.lg.WriteLog("StreamServerConn", err.Error())
					closed <- struct{}{}
				}
			}
		case <-done:
			total--
			if total == 0 {
				close(done)
				close(closed)
				time.Sleep(5 * time.Millisecond)
				close(stream)
				close(errMsg)
				return fmt.Errorf("Connection closed")
			}
		}
	}
}

func (c *grpcService) StreamBiConn(stream repos.StreamBiInt) error {
	closedTx := make(chan struct{}, 0)
	closedRx := make(chan struct{}, 0)
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

	go c.client.BidirectionalStream(subscribers, streamProv, streamSub, errMsg, closedTx, closedRx, done)
	streamProv <- req.GetPayload()

	//Receives from provider
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				go c.lg.WriteLog("StreamBiConn", err.Error())
				close(streamProv)
				<-closedTx
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
					closedRx <- struct{}{}
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

type GrpcServiceInt interface {
	SimpleConn(ms repos.GrpRequestMessageInt) error
	StreamClientConn(st repos.StreamClientInt) error
	StreamServerConn(req repos.GrpRequestMessageInt, st repos.StreamServerInt) error
	StreamBiConn(st repos.StreamBiInt) error
}
