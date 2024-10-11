package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"dominus/app/domain/topic"
	"dominus/app/interfaces/database"
	"dominus/app/interfaces/rest/output"
	"sync"
	"time"
)

type connection struct {
	m *entities.Message // Message
	// grpc struct
	t    topic.TopicInt         // Topic
	ev   event.EventsInt        // Events
	r    rules.RuleInt          // Rules
	cr   output.RestClientInt   // Rest client
	repo database.RepositoryInt // Repository client
	// grpc client
	lg event.LogsInt
}

func (c *connection) RestToGrpc(ctx RestContextInt) error {


	return nil
}

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

func (c *connection) GrpcToRest() error {

	
	return nil
}

func (c *connection) GrpcToGrpc() error {
	return nil
}

type ConnectionInt interface {
	//RestToGrp sends information from  rest protocol receiver to  grpc protocol client
	//
	//Parameters
	//
	//-> ctx: RestContexInt
	RestToGrpc(ctx RestContextInt) error
	//RestToRest sends information from rest protocol receiver to rest protocol client
	//
	//Parameters
	//
	//-> ctx: RestContexInt
	RestToRest(ctx RestContextInt) error
	//GrpcToRest sends information from grpc protocol receiver to rest protocol client
	GrpcToRest() error
}
