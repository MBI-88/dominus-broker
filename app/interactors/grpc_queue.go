package interactors

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/domain/repos"
	"fmt"
	"time"
)

func (c *grpcService) SimpleConn(ms repos.GrpRequestMessageInt) error {
	name := ms.GetSubscribers()
	if len(name) == 0 {
		return fmt.Errorf("Topic name empty")
	}
	topic, err := c.topics.Find(name[0])
	if err != nil {
		return err
	}
	if len(topic.Partitions) == 0 {
		return fmt.Errorf("Partitions are empty")
	}

	return topic.Push(ms.GetPayload())
}

func (c *grpcService) checkQueue() {
	for {
		t, err := c.topics.Next()
		if err != nil {
			break
		}

		for _, sub := range t.Partitions {
			go func(url string, topic entities.Topic) {
				body := topic.Pop()
				resp, err := c.client.Simple(url, body)
				if err != nil {
					c.lg.WriteLog("SimpleConn", err.Error())
					c.glue(t.Name, body)
					return
				}
				if resp.GetStatus() != 200 {
					c.glue(t.Name, body)
				}
			}(sub, t)
		}
	}
}

func (c *grpcService) glue(name string, msg []byte) {
	topic, err := c.topics.Find(name)
	if err != nil {
		c.lg.WriteLog("glue", err.Error())
		return
	}
	if err := topic.Push(msg); err != nil {
		c.lg.WriteLog("glue", err.Error())
	}
}

func (c *grpcService) RunQueue() {
	for {
		select {
		case <-time.Tick(2 * time.Second):
			c.checkQueue()
		}
	}
}
