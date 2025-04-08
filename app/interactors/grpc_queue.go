package interactors

import (
	"dominus-project/app/domain/repos"
	"fmt"
)

func (c *grpcService) SimpleConn(ms repos.GrpRequestMessageInt) error {
	subs := ms.GetSubscribers()
	if len(subs) == 0 {
		return fmt.Errorf("Subscribers not found")
	}
	body := ms.GetPayload()
	for _, sub := range subs {
		go func(url string, body []byte) {
			_, err := c.client.Simple(url, body)
			if err != nil {
				c.lg.WriteLog("SimpleConn", err.Error())
			}
		}(sub, body)
	}
	return nil
}
