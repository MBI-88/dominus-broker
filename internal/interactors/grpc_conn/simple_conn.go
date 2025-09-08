package grpcconn

import (
	"dominus-project/internal/domain/adapters"
	"fmt"
)

func (c *grpcService) SimpleConn(ms adapters.GrpcDto) error {
	name := c.getTopicName(ms)
	if name == "" {
		return fmt.Errorf("topic name empty")
	}
	topic, err := c.topics.Find(name)
	if err != nil {
		return err
	}
	if len(topic.GetSubscribers()) == 0 {
		return fmt.Errorf("subscribers are empty")
	}

	message := ms.GetPayload()
	if len(message) == 0 {
		return fmt.Errorf("message empty")
	}

	return topic.SetMessage(message)
}
