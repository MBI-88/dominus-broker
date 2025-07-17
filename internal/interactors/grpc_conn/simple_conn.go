package grpcconn

import (
	"dominus-project/internal/domain/adapters"
	"fmt"
)

func (c *grpcService) SimpleConn(ms adapters.IGrpcDto) error {
	name := c.getTopicName(ms)
	if name == "" {
		return fmt.Errorf("Topic name empty")
	}
	topic, err := c.topics.Find(name)
	if err != nil {
		return err
	}
	if len(topic.GetSubscribers()) == 0 {
		return fmt.Errorf("Subscribers are empty")
	}

	message := ms.GetPayload()
	if len(message) == 0 {
		return fmt.Errorf("Message empty")
	}

	return topic.SetMessage(message)
}
