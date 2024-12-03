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
	opts []grpc.DialOption
	rls  rules.RulesInt
}

func (g *grpcClient) Simple(url string, body []byte) (interactors.GrpResponseInt, error) {
	conn, err := grpc.NewClient(url, g.opts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewGrpcClient(conn)
	msg := &pb.RequestMessage{
		Subscribers: nil, Payload: body,
	}

	resp, err := client.Simple(context.Background(), msg)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *grpcClient) ClientStream(urls []string, msg <-chan []byte, sig chan<- entities.Logs) {
	arrayMsg := make([]chan []byte, 0, len(urls))
	for p, url := range urls {
		ch := make(chan []byte, 0)
		arrayMsg = append(arrayMsg, ch)
		go func(url string, ch chan []byte, p int) {
			if ok := g.rls.CheckURI(url); ok {
				conn, err := grpc.NewClient(url, g.opts...)
				if err == nil {
					client := pb.NewGrpcClient(conn)
					stream, err := client.ClientStream(context.Background())
					if err == nil {
						go func(client pb.Grpc_ClientStreamClient, p int) {
							for {
								select {
								case payload, ok := <-ch:
									if ok {
										if err := client.Send(&pb.RequestMessage{
											Subscribers: urls,
											Payload:     payload}); err != nil {
											sig <- entities.Logs{
												CreatedAt:  time.Now(),
												Desc:       err.Error(),
												Stage:      "ClientStream sends to subscribers",
												Subscriber: urls[p],
											}
										}
									} else {
										client.CloseSend()
										return
									}

								}
							}
						}(stream, p)
					}
				}
			}
		}(url, ch, p)
	}

	for {
		select {
		case payload, ok := <-msg:
			if ok {
				for _, ch := range arrayMsg {
					ch <- payload
				}
			} else {
				for _, ch := range arrayMsg {
					close(ch)
				}
				return
			}
		}
	}
}

func (g *grpcClient) ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, sig chan<- entities.Logs, ctx context.Context) {
	for p, url := range urls {
		go func(url string, p int) {
			if ok := g.rls.CheckURI(url); ok {
				conn, err := grpc.NewClient(url, g.opts...)
				if err == nil {
					client := pb.NewGrpcClient(conn)
					reqMsg := &pb.RequestMessage{
						Subscribers: nil,
						Payload:     initalMsg,
					}
					stream, err := client.ServerStream(context.Background(), reqMsg)
					if err == nil {
						go func(client pb.Grpc_ServerStreamClient, p int) {
							for {
								resp, err := client.Recv()
								if err != nil {
									sig <- entities.Logs{
										Desc:       err.Error(),
										CreatedAt:  time.Now(),
										Subscriber: urls[p],
										Stage:      "ServerStream receives from subscribers",
									}
									return
								} else {
									if _, ok := <-ctx.Done(); !ok {
										msg <- resp.Payload
									} else {
										return
									}
								}
							}
						}(stream, p)
					}
				}
			}
		}(url, p)
	}
}

func (g *grpcClient) BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- entities.Logs, ctx context.Context) {
	arrayMsg := make([]chan []byte, 0, len(urls))
	for p, url := range urls {
		ch := make(chan []byte, 0)
		arrayMsg = append(arrayMsg, ch)
		go func(ch <-chan []byte, p int, sub string) {
			if ok := g.rls.CheckURI(url); ok {
				conn, err := grpc.NewClient(url, g.opts...)
				if err == nil {
					client := pb.NewGrpcClient(conn)
					stream, err := client.BidirectionalStream(context.Background())
					if err == nil {
						//Sends to subscriber
						go func(client pb.Grpc_BidirectionalStreamClient, p int) {
							for {
								select {
								case payload, ok := <-ch:
									if ok {
										if err := client.Send(&pb.RequestMessage{
											Subscribers: urls,
											Payload:     payload}); err != nil {
											errMsg <- entities.Logs{
												CreatedAt:  time.Now(),
												Desc:       err.Error(),
												Stage:      "BidirectionalStream sends to subscribers",
												Subscriber: urls[p],
											}
										}
									} else {
										client.CloseSend()
										return
									}
								}
							}
						}(stream, p)

						//Receives from subscriber
						go func(client pb.Grpc_BidirectionalStreamClient, p int) {
							for {
								resp, err := client.Recv()
								if err != nil {
									log := entities.Logs{
										CreatedAt:  time.Now(),
										Stage:      "BidirectionalStream Recv from subscribers",
										Subscriber: urls[p],
										Desc:       err.Error(),
									}
									errMsg <- log
									return
								}
								if _, ok := <-ctx.Done(); !ok {
									subMsg <- resp.GetPayload()
								} else {
									return
								}
							}
						}(stream, p)
					}
				}
			}
		}(ch, p, url)
	}

	for {
		select {
		case payload, ok := <-provMsg:
			if ok {
				for _, ch := range arrayMsg {
					ch <- payload
				}
			} else {
				for _, ch := range arrayMsg {
					close(ch)
				}
				return
			}
		}
	}
}

func NewGrpClient(opts []grpc.DialOption) interactors.GrpClientInt {
	return &grpcClient{
		opts: opts,
		rls:  rules.NewRule(),
	}
}
