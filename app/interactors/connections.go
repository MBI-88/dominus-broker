package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"fmt"
	"time"
)

type connection struct {
	r      rules.RuleInt // Rules
	cr     RestClientInt // Rest client
	repo   RepositoryInt // Repository client
	client GrpClientInt
	lg     event.LogsInt
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
				resp, err := c.client.Simple(url, body)
				if err != nil {
					logs := entities.Logs{
						Desc:      err.Error(),
						CreatedAt: time.Now(),
						Status:    resp.GetStatus(),
						Stage:     "SimpleConn",
					}

					if _, err := c.repo.InsertObject(&logs, "logs"); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
				}

			}(sub, body)
		}
	}
	return nil
}

func (c *connection) StreamClientConn(st StreamClientInt) error {
	stream := make(chan []byte)
	errMsg := make(chan *entities.Logs)
	req, err := st.Recv()
	if err != nil {
		return err
	}

	go c.client.ClientStream(req.GetSubscribers(), stream, errMsg)
	go func(sig <-chan *entities.Logs) {
	loop:
		for {
			select {
			case val, ok := <-sig:
				if ok {
					if _, err := c.repo.InsertObject(val, "logs"); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
				} else {
					break loop
				}
			}
		}
	}(errMsg)

	stream <- req.GetPayload()
	for {
		req, err := st.Recv()
		if err != nil {
			close(stream)
			close(errMsg)
			return err
		}
		stream <- req.GetPayload()
	}
}

func (c *connection) StreamServerConn(req GrpRequestMessageInt, st StreamServerInt) error {
	stream := make(chan []byte, len(req.GetSubscribers()))
	errMsg := make(chan *entities.Logs)
	initialRequest := req.GetPayload()

	go c.client.ServerStream(req.GetSubscribers(), initialRequest, stream, errMsg)

	go func(sig <-chan *entities.Logs) {
	loop:
		for {
			select {
			case val, ok := <-sig:
				if ok {
					if _, err := c.repo.InsertObject(val, "logs"); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
				} else {
					break loop
				}
			}
		}
	}(errMsg)

loop:
	for {
		select {
		case body, ok := <-stream:
			if ok {
				if err := st.Send(body); err != nil {
					log := &entities.Logs{
						Desc:      err.Error(),
						CreatedAt: time.Now(),
						Stage:     "StreamServerConn Send to provider",
						Status:    uint32(500),
					}
					if _, err := c.repo.InsertObject(log, "logs"); err != nil {
						c.lg.WriteLog("InsertObject", err.Error())
					}
					close(stream)
					close(errMsg)
				}
			} else {
				break loop
			}

		}
	}
	return fmt.Errorf("Connection closed")
}

func (c *connection) StreamBiConn(stream StreamBiInt) error {

	return nil
}

type ConnectionInt interface {
	SimpleConn(ms GrpRequestMessageInt) error
	StreamClientConn(st StreamClientInt) error
	StreamServerConn(req GrpRequestMessageInt, st StreamServerInt) error
	StreamBiConn(st StreamBiInt) error
}
