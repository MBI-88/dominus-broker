package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/repos"
	"dominus/app/domain/rules"
	"fmt"
	"time"
)

type connection struct {
	r          rules.RulesInt // Rules
	repo       repos.RepositoryInt  // Repository client
	client     repos.GrpClientInt
	lg         event.LogsInt
	collection string
}

func (c *connection) SimpleConn(ms repos.GrpRequestMessageInt) error {
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

func (c *connection) StreamClientConn(st repos.StreamClientInt) error {
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
	errMsg := make(chan entities.Logs, len(subscribers))

	go c.client.ClientStream(subscribers, stream, errMsg, closed)
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
			log := entities.Logs{
				ID:        c.r.MakeID(),
				Desc:      err.Error(),
				CreatedAt: time.Now(),
				Stage:     "StreamClientConn receives from provider",
			}
			if err := c.repo.InsertObject(log, c.collection); err != nil {
				c.lg.WriteLog("InsertObject", err.Error())
			}
			close(stream)
			<-closed
			close(errMsg)
			return err
		}
		stream <- req.GetPayload()
	}
}

func (c *connection) StreamServerConn(req repos.GrpRequestMessageInt, st repos.StreamServerInt) error {
	closed := make(chan struct{}, 0)
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		return fmt.Errorf("Subscribers not found")
	}
	total := len(subscribers)
	stream := make(chan []byte, total+int(total*2/3))
	errMsg := make(chan entities.Logs, total)
	done := make(chan struct{}, total)
	initialRequest := req.GetPayload()

	go c.client.ServerStream(subscribers, initialRequest, stream, errMsg, closed, done)

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

func (c *connection) StreamBiConn(stream repos.StreamBiInt) error {
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
	errMsg := make(chan entities.Logs, total)

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

	go c.client.BidirectionalStream(subscribers, streamProv, streamSub, errMsg, closedTx, closedRx, done)
	streamProv <- req.GetPayload()

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
					log := entities.Logs{
						ID:        c.r.MakeID(),
						Desc:      err.Error(),
						CreatedAt: time.Now(),
						Stage:     "StreamBiConn Send to provider",
					}
					if err := c.repo.InsertObject(log, c.collection); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
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

type ConnectionInt interface {
	SimpleConn(ms repos.GrpRequestMessageInt) error
	StreamClientConn(st repos.StreamClientInt) error
	StreamServerConn(req repos.GrpRequestMessageInt, st repos.StreamServerInt) error
	StreamBiConn(st repos.StreamBiInt) error
}
