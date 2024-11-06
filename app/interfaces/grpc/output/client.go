package output

import (
	"context"
	"dominus/app/domain/entities"
	"dominus/app/domain/rules"
	"dominus/app/interactors"
	pb "dominus/app/interfaces/grpc/proto/builder"
	"time"

	"google.golang.org/grpc"
)

type grpcClient struct {
	opts         []grpc.DialOption
	rls          rules.RuleInt
	clientStream []pb.Grpc_ClientStreamClient
	serverStream []pb.Grpc_ServerStreamClient
}

func (g *grpcClient) Simple(url string, body []byte) (interactors.GrpResponseInt, error) {
	conn, err := grpc.NewClient(url, g.opts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewGrpcClient(conn)
	ctx := context.Background()

	msg := &pb.RequestMessage{
		Subscribers: []string{}, Payload: body,
	}

	resp, err := client.Simple(ctx, msg)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *grpcClient) ClientStream(urls []string, msg <-chan []byte, sig chan<- *entities.Logs) {
	for _, url := range urls {
		if ok := g.rls.CheckURI(url); ok {
			conn, err := grpc.NewClient(url, g.opts...)
			if err == nil {
				client := pb.NewGrpcClient(conn)
				stream, err := client.ClientStream(context.Background())
				if err == nil {
					g.clientStream = append(g.clientStream, stream)
				}
			}
		}
	}

	if len(g.clientStream) > 0 {
	loop:
		for {
			select {
			case payload, ok := <-msg:
				if ok {
					for i, c := range g.clientStream {
						go func(client pb.Grpc_ClientStreamClient, cnumber int) {
							if err := client.Send(&pb.RequestMessage{
								Subscribers: urls,
								Payload:     payload}); err != nil {

								sig <- &entities.Logs{
									CreatedAt:  time.Now(),
									Desc:       err.Error(),
									Status:     uint32(500),
									Stage:      "ClientStream send to subscribers",
									Subscriber: urls[i],
								}
							}
						}(c, i)
					}

				} else {
					break loop
				}

			}
		}
	}

}

func (g *grpcClient) ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, sig chan<- *entities.Logs) {
	for _, url := range urls {
		if ok := g.rls.CheckURI(url); ok {
			conn, err := grpc.NewClient(url, g.opts...)
			if err == nil {
				client := pb.NewGrpcClient(conn)
				reqMsg := &pb.RequestMessage{
					Subscribers: urls,
					Payload:     initalMsg,
				}
				stream, err := client.ServerStream(context.Background(), reqMsg)
				if err == nil {
					g.serverStream = append(g.serverStream, stream)
				}
			}
		}
	}

	if len(g.serverStream) > 0 {

	}

}

func (g *grpcClient) BidirectionalStream(url string) {

}

func NewGrpClient(opts []grpc.DialOption) interactors.GrpClientInt {
	return &grpcClient{
		opts: opts,
		rls:  rules.NewRule(),
	}
}
