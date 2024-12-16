package interactors

import (
	"context"
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"fmt"
	"time"
)

type connection struct {
	r          rules.RulesInt // Rules
	repo       RepositoryInt  // Repository client
	client     GrpClientInt
	lg         event.LogsInt
	collection string
}

func (c *connection) SimpleConn(ms GrpRequestMessageInt) error {
	subs := ms.GetSubscribers()
	if len(subs) == 0 {
		return fmt.Errorf("Subscribers not found")
	}
	body := ms.GetPayload()
	for _, sub := range subs {
		if ok := c.r.CheckURI(sub); ok {
			go func(url string, body []byte) {
				_, err := c.client.Simple(url, body)
				if err != nil {
					logs := entities.Logs{
						ID:        c.r.MakeID(),
						Desc:      err.Error(),
						CreatedAt: time.Now(),
						Stage:     "SimpleConn",
					}
					if err := c.repo.InsertObject(logs, c.collection); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
				}
			}(sub, body)
		}
	}
	return nil
}

func (c *connection) StreamClientConn(st StreamClientInt) error {
	stream := make(chan []byte, 0)
	errMsg := make(chan entities.Logs, 0)
	isWatting := make(chan struct{}, 0)
	req, err := st.Recv()
	if err != nil {
		return err
	}

	go c.client.ClientStream(req.GetSubscribers(), stream, errMsg, isWatting)
	go func(sig <-chan entities.Logs) {
		for {
			select {
			case val, ok := <-sig:
				if ok {
					if err := c.repo.InsertObject(val, c.collection); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
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
			close(stream)
			<-isWatting
			close(errMsg)
			return err
		}
		stream <- req.GetPayload()
	}
}

func (c *connection) StreamServerConn(req GrpRequestMessageInt, st StreamServerInt) error {
	counter := 0
	cal := make(chan struct{}, 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream := make(chan []byte, len(req.GetSubscribers()))
	errMsg := make(chan entities.Logs, 0)
	initialRequest := req.GetPayload()

	go c.client.ServerStream(req.GetSubscribers(), initialRequest, stream, errMsg, ctx)

	go func(sig <-chan entities.Logs) {
		for {
			select {
			case val, ok := <-sig:
				if ok {
					if err := c.repo.InsertObject(val, c.collection); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
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
					log := entities.Logs{
						ID:        c.r.MakeID(),
						Desc:      err.Error(),
						CreatedAt: time.Now(),
						Stage:     "StreamServerConn Send to provider",
					}
					if err := c.repo.InsertObject(log, c.collection); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
					cancel()
					close(stream)
					close(errMsg)
				}
				counter = 0
			} else {
				close(cal)
				return fmt.Errorf("Connection closed")
			}
		case <-time.Tick(3 * time.Second):
			counter++
			if counter >= 2 {
				cal <- struct{}{}
			}
		case <-cal:
			cancel()
			close(stream)
			close(errMsg)

		}
	}
}

func (c *connection) StreamBiConn(stream StreamBiInt) error {
	counter := 0
	cal := make(chan struct{}, 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	request, err := stream.Recv()
	subscribers := request.GetSubscribers()
	streamProv := make(chan []byte, 0)
	streamSub := make(chan []byte, len(subscribers))
	errMsg := make(chan entities.Logs, 0)

	if err != nil {
		return err
	}
	go func(sig <-chan entities.Logs) {
		for {
			select {
			case val, ok := <-sig:
				if ok {
					if err := c.repo.InsertObject(val, c.collection); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
				} else {
					return
				}
			}
		}
	}(errMsg)

	go c.client.BidirectionalStream(subscribers, streamProv, streamSub, errMsg, ctx)
	streamProv <- request.GetPayload()

	//Receives from provider
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				log := entities.Logs{
					ID:        c.r.MakeID(),
					Desc:      err.Error(),
					CreatedAt: time.Now(),
					Stage:     "StreamBiConn Recv from provider",
				}
				if err := c.repo.InsertObject(log, c.collection); err != nil {
					c.lg.WriteLog("InsertObject", err.Error())
				}
				cancel()
				close(streamProv)
				close(streamSub)
				close(errMsg)
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
					log := entities.Logs{
						ID:        c.r.MakeID(),
						Desc:      err.Error(),
						CreatedAt: time.Now(),
						Stage:     "StreamBiConn Send to provider",
					}
					if err := c.repo.InsertObject(log, c.collection); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
					cancel()
					close(streamProv)
					close(streamSub)
					close(errMsg)
				}
				counter = 0
			} else {
				close(cal)
				return fmt.Errorf("Connection closed")
			}
		case <-time.Tick(3 * time.Second):
			counter++
			if counter >= 2 {
				cal <- struct{}{}
			}
		case <-cal:
			cancel()
			close(streamProv)
			close(streamSub)
			close(errMsg)
		}
	}
}

type ConnectionInt interface {
	SimpleConn(ms GrpRequestMessageInt) error
	StreamClientConn(st StreamClientInt) error
	StreamServerConn(req GrpRequestMessageInt, st StreamServerInt) error
	StreamBiConn(st StreamBiInt) error
}
