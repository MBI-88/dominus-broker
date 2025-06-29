package grpcconn

import (
	"dominus-project/app/domain/repos"
	"fmt"
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
	if len(topic.Subscribers) == 0 {
		return fmt.Errorf("Partitions are empty")
	}

	return topic.Push(ms.GetPayload())
}



