package broker

import (
	"context"

	"dominus-project/internal/application/dtos"
	"dominus-project/internal/domain/enum"
	"fmt"
)

func (b *broker) StreamClientConn(st dtos.BrokerClientDto) error {
	ctx, cancel := context.WithCancel(st.Context())
	defer cancel()
	stream := make(chan []byte)
	req, err := st.Recv()
	if err != nil {
		go b.lg.WriteLog(ctx, enum.ERROR, "StreamClientConn", err.Error())
		return err
	}
	subscribers := req.GetSubscribers()
	if len(subscribers) == 0 {
		go b.lg.WriteLog(ctx, enum.ERROR, "StreamClientConn", enum.SUBCRIBER_NOT_FOUND)
		return fmt.Errorf("subscribers not found")
	}

	go b.client.ClientStream(subscribers, stream, ctx)

	stream <- req.GetPayload()
	for {
		req, err := st.Recv()
		if err != nil {
			go b.lg.WriteLog(ctx, enum.ERROR ,"StreamClientConn", err.Error())
			close(stream)
			cancel()
			return err
		}
		stream <- req.GetPayload()
	}
}
