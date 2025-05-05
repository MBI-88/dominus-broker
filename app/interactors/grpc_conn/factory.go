package grpcconn

import (
	"dominus-project/app/domain/entities"
	"dominus-project/app/domain/repos"
)

type grpcService struct {
	client repos.GrpClientInt
	lg     repos.LogsInt
	topics entities.TopicsInt
}


func (c *grpcService) checkQueue() {
	for {
		t, err := c.topics.Next()
		if err != nil {
			break
		}

		for _, sub := range t.Subscribers {
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

type GrpcServiceInt interface {
	SimpleConn(ms repos.GrpRequestMessageInt) error
	StreamClientConn(st repos.StreamClientInt) error
	StreamServerConn(req repos.GrpRequestMessageInt, st repos.StreamServerInt) error
	StreamBiConn(st repos.StreamBiInt) error
	RunQueue()
}


func NewGrpcService( lclient repos.LogsInt, gclient repos.GrpClientInt, topics entities.TopicsInt) GrpcServiceInt {
	return &grpcService{
		lg:     lclient,
		client: gclient,
		topics: topics,
	}
}