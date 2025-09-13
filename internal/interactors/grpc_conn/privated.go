package grpcconn

import (
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/entities"
)

// CheckQueue Check the Queue for new messages
//
// # Params
//
// Empty
func (c *grpcService) checkQueue() {
	for {
		t, err := c.topics.Next()
		if err != nil {
			c.lg.WriteLog("Next", err.Error())
			break
		}

		body := t.GetMessage()
		if len(body) > 0 {
			for _, sub := range t.GetSubscribers() {
				go func(url string, topic entities.Topic) {
					resp, err := c.client.Simple(url, body)
					if err != nil {
						c.lg.WriteLog("SimpleConn", err.Error())
						if err = topic.SetMessage(body); err != nil {
							c.lg.WriteLog("SetMessage", err.Error())
						}
						return
					}
					if resp.GetStatus() != 200 {
						if err = topic.SetMessage(body); err != nil {
							c.lg.WriteLog("SetMessage", err.Error())
						}
					}
				}(sub, t)
			}
		}

	}
}

func (*grpcService) getTopicName(ms adapters.GrpcDto) string {
	if subs := ms.GetSubscribers(); len(subs) > 0 {
		return subs[0]
	}
	return ""
}
