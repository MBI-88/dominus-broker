package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"dominus/app/domain/topic"
	"fmt"
	"time"
)

type connection struct {
	t      topic.TopicInt  // Topic
	ev     event.EventsInt // Events
	r      rules.RuleInt   // Rules
	cr     RestClientInt   // Rest client
	repo   RepositoryInt   // Repository client
	client GrpClientInt
	lg     event.LogsInt
}

func (c *connection) SimpleConn(ms GrpRequestMessageInt) error {
	topic := ms.GetTopic()

	if topic == "" {
		return fmt.Errorf("Not topic found")
	}

	subs := c.t.GetSusbcribers(topic)
	if len(subs) == 0 {
		return fmt.Errorf("Subscribers emtpy")
	}

	body := ms.GetPayload()
	for _, sub := range subs {
		go func(url string, body []byte) {
			resp, err := c.client.Simple(url, body)
			if err != nil {
				logs := entities.Logs{
					Desc:      err.Error(),
					CreatedAt: time.Now(),
					Status:    resp.GetStatus(),
					Sub:       sub,
				}

				if _, err := c.repo.InsertObject(&logs, "logs"); err != nil {
					c.lg.WriteLog("InsertObject", err.Error())
				}
			}

		}(sub, body)
	}

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
