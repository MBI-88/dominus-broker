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

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type connection struct {
	m *entities.Message // Message
	// grpc struct
	p    jsoniter.API           // Parser
	t    topic.TopicInt         // Topic
	ev   event.EventsInt        // Events
	r    rules.RuleInt          // Rules
	cr   output.RestClientInt   // Rest client
	repo database.RepositoryInt // Repository client
	// grpc client
}

func (c connection) RestToGrpc(ctx *fasthttp.RequestCtx) {

}

func (c connection) RestToRest(ctx *fasthttp.RequestCtx) {
	var (
		message = make(map[string]any)
		body    = ctx.Request.Body()
	)

	if err := c.p.Unmarshal(body, c.m); err != nil {
		message["message"] = err.Error()
		b, _ := c.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	if err := c.r.ValidateStruct(c.m); err != nil {
		message["message"] = err.Error()
		b, _ := c.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	subs := c.t.GetSusbcribers(c.m.Topic)
	ch := make(chan bool, len(subs))
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
							// Create a client log
						}
					}(log)

					sig <- true
				}
			}(sub, msg, &wg )

		}
	
		wg.Wait()
		close(sig)
	}(subs, c.m, ch)

	message["message"] = "Accepted"
	b, _ := c.p.Marshal(message)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(b)
}

func (c connection) GrpcToRest() {

}

func (c connection) GrpcToGrpc() {

}

type ConnectionInt interface {
	//RestToGrp sends information from  rest protocol receiver to  grpc protocol client
	//
	//Parameters
	//
	//-> ctx: fasthttp context
	RestToGrpc(ctx *fasthttp.RequestCtx)
	//RestToRest sends information from rest protocol receiver to rest protocol client
	//
	//Parameters
	//
	//-> ctx: fasthttp context
	RestToRest(ctx *fasthttp.RequestCtx)
	//GrpcToRest sends information from grpc protocol receiver to rest protocol client
	GrpcToRest()
}
