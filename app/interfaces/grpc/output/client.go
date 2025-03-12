package output

import (
	"context"
	"github.com/PR0C0D3-MBI/dominus-project/app/domain/repos"
	"github.com/PR0C0D3-MBI/dominus-project/app/domain/rules"
	pb "github.com/PR0C0D3-MBI/dominus-project/app/interfaces/grpc/proto/builder"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

type grpcClient struct {
	opts []grpc.DialOption
	rls  rules.RulesInt
}

func (g *grpcClient) Simple(url string, body []byte) (repos.GrpResponseInt, error) {
	if g.rls.CheckURI(url) {
		conn, err := grpc.NewClient(url, g.opts...)
		if err != nil {
			return nil, err
		}
		client := pb.NewGrpcClient(conn)
		msg := &pb.RequestMessage{
			Subscribers: nil, Payload: body,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		resp, err := client.Simple(ctx, msg)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, fmt.Errorf("Invalid uri")
}

func (g *grpcClient) ClientStream(urls []string, msg <-chan []byte, sig chan<- error, tx chan<- struct{}) {
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
							isClose := false
							for {
								select {
								case payload, ok := <-ch:
									if ok {
										if !isClose {
											if err := client.Send(&pb.RequestMessage{
												Subscribers: urls,
												Payload:     payload}); err != nil {
												sig <- err
												isClose = true
											}
										}
									} else {
										return
									}
								}
							}
						}(stream, p)
					} else {
						return
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
				tx <- struct{}{}
				return
			}
		}
	}
}

func (g *grpcClient) ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, sig chan<- error, ctx context.Context, tx chan<- struct{}) {
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
					stream, err := client.ServerStream(ctx, reqMsg)
					if err == nil {
						go func(c pb.Grpc_ServerStreamClient, p int) {
							for {
								resp, err := c.Recv()
								if err != nil {
									sig <- err
									tx <- struct{}{}
									return
								} else {
									msg <- resp.GetPayload()
								}
							}
						}(stream, p)
					} else {
						return
					}
				}
			}
		}(url, p)
	}
}

func (g *grpcClient) BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- error, tx chan<- struct{}, ctx context.Context, done chan<- struct{}) {
	arrayMsg := make([]chan []byte, 0, len(urls))
	for p, url := range urls {
		ch := make(chan []byte, 0)
		arrayMsg = append(arrayMsg, ch)
		go func(ch <-chan []byte, p int, sub string) {
			if ok := g.rls.CheckURI(url); ok {
				conn, err := grpc.NewClient(url, g.opts...)
				if err == nil {
					client := pb.NewGrpcClient(conn)
					stream, err := client.BidirectionalStream(ctx)
					if err == nil {
						//Sends to subscriber
						go func(c pb.Grpc_BidirectionalStreamClient, p int) {
							for {
								select {
								case payload, ok := <-ch:
									if ok {
										if err := c.Send(&pb.RequestMessage{
											Subscribers: urls,
											Payload:     payload}); err != nil {
											errMsg <- err
										}
									} else {
										return
									}
								}
							}
						}(stream, p)

						//Receives from subscriber
						go func(c pb.Grpc_BidirectionalStreamClient, p int) {
							for {
								resp, err := c.Recv()
								if err != nil {
									errMsg <- err
									done <- struct{}{}
									return
								}
								subMsg <- resp.GetPayload()
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
				tx <- struct{}{}
				return
			}
		}
	}
}

func NewGrpClient(opts []grpc.DialOption) repos.GrpClientInt {
	return &grpcClient{
		opts: opts,
		rls:  rules.NewRules(),
	}
}
