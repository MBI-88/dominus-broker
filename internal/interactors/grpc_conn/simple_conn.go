package grpcconn

import (
	"dominus-project/internal/domain/repos"
	"fmt"
)

func (c *grpcService) SimpleConn(ms repos.GrpRequestMessageInt) error {
	name := c.getTopicName(ms)
	if name == "" {
		return fmt.Errorf("Topic name empty")
	}
	topic, err := c.topics.Find(name)
	if err != nil {
		return err
	}
	if len(topic.Subscribers) == 0 {
		return fmt.Errorf("Subscribers are empty")
	}
	
	message := ms.GetPayload()
	if len(message) == 0 {
		return  fmt.Errorf("Message empty")
	}
	topic.SetMessage(message)
	return nil
}
