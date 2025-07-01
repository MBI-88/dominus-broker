package grpcconn

import "dominus-project/internal/domain/entities"

// CheckQueue Check the Queue for new messages
//
// Params
//
// Empty
func (c *grpcService) checkQueue() {
	for {
		t, err := c.topics.Next()
		if err != nil {
			break
		}

		for _, sub := range t.Subscribers {
			go func(url string, topic *entities.Topic) {
				body := topic.Pop()
				resp, err := c.client.Simple(url, body)
				if err != nil {
					c.lg.WriteLog("SimpleConn", err.Error())
					topic.Push(body)
					return
				}
				if resp.GetStatus() != 200 {
					topic.Push(body)
				}
			}(sub, t)
		}
	}
}