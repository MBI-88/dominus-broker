package grpcconn

import (
	"dominus-project/internal/domain/entities"
	"dominus-project/internal/domain/repos"
)

// CheckQueue Check the Queue for new messages
//
// Params
//
// Empty
//
func (c *grpcService) checkQueue() {
	for {
		t, err := c.topics.Next()
		if err != nil {
			break
		}

		for _, sub := range t.GetSubscribers() {
			go func(url string, topic entities.TopicInt) {
				body := topic.GetMessage()
				if len(body) > 0 {
					resp, err := c.client.Simple(url, body)
					if err != nil {
						c.lg.WriteLog("SimpleConn", err.Error())
						topic.SetMessage(body)
						return
					}
					if resp.GetStatus() != 200 {
						topic.SetMessage(body)
					}
				}

			}(sub, t)
		}
	}
}

func (*grpcService) getTopicName(ms repos.GrpRequestMessageInt) string {
	if len(ms.GetSubscribers()) > 0 {
		return  ms.GetSubscribers()[0]
	}
	return  ""
}