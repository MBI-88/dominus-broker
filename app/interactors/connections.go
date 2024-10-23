package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"dominus/app/domain/topic"
	"sync"
	"time"
)

type connection struct {
	m *entities.Message // Message
	// grpc struct
	t    topic.TopicInt  // Topic
	ev   event.EventsInt // Events
	r    rules.RuleInt   // Rules
	cr   RestClientInt   // Rest client
	repo RepositoryInt   // Repository client
	// grpc client
	lg event.LogsInt
}


// Firt iteration

func (c *connection) RestToRest(ctx RestContextInt) error {
	if err := ctx.BodyParser(c.m); err != nil {
		c.lg.WriteLog("InsertObject", err.Error())
		return err
	}

	if err := c.r.ValidateStruct(c.m); err != nil {
		c.lg.WriteLog("InsertObject", err.Error())
		return err
	}

	subs := c.t.GetSusbcribers(c.m.Topic)
	ch := make(chan bool)
	go c.ev.Sentinel(ch)

	go func(subscribers []string, msg *entities.Message, sig chan<- bool) {
		var wg sync.WaitGroup

		for _, sub := range subscribers {
			wg.Add(1)
			go func(addr string, msg *entities.Message, wg *sync.WaitGroup) {
				defer wg.Done()
				if err := c.cr.DoJsonRequest(sub, msg.Payload); err != nil {
					log := &entities.Logs{
						Log:       *msg,
						CreatedAt: time.Now(),
						Status:    "pending",
					}

					go func(log *entities.Logs) {
						if _, err := c.repo.InsertObject(log, "watting for definition"); err != nil {
							c.lg.WriteLog("InsertObject", err.Error())
						}
					}(log)

					sig <- true
				}
			}(sub, msg, &wg)

		}

		wg.Wait()
		close(sig)
	}(subs, c.m, ch)

	return nil
}



func (c *connection) SimpleConn(ms GrpRequestMessageInt) error {

	return nil
}

func (c *connection) StreamClientConn(stream StreamClientInt) error {

	return nil
}

func (c *connection) StreamServerConn(req GrpRequestMessageInt, stream StreamServerInt) error {
	
	return nil
}

func (c *connection) StreamBiConn(stream StreamBiInt) error {

	return nil
}



type ConnectionInt interface {
	SimpleConn(ms GrpRequestMessageInt) error
	StreamClientConn(stream StreamClientInt) error
	StreamServerConn(req GrpRequestMessageInt, stream StreamServerInt) error
	StreamBiConn(stream StreamBiInt) error
}
