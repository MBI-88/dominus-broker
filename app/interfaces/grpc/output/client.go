package output

import (
	"context"
	"dominus-project/app/domain/repos"
	pb "dominus-project/app/interfaces/grpc/proto/builder"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type grpcClient struct {
	opts []grpc.DialOption
	lgs  repos.LogsInt
}

func (g *grpcClient) Simple(url string, body []byte) (repos.GrpResponseInt, error) {
	conn, err := grpc.NewClient(url, g.opts...)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
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

func (g *grpcClient) ClientStream(urls []string, msg <-chan []byte) {
	arrayDoQuery := make([]func([]byte), 0, len(urls))
	lock := new(sync.Mutex)
	connect := func(url string) (pb.Grpc_ClientStreamClient, error) {
		conn, _ := grpc.NewClient(url, g.opts...)
		c := pb.NewGrpcClient(conn)
		return c.ClientStream(context.Background())
	}
	doQuery := func(url string) func([]byte) {
		stream, errConn := connect(url)
		return func(payload []byte) {
			lock.Lock()
			err := errConn
			lock.Unlock()
			if err == nil {
				if err := stream.Send(&pb.RequestMessage{
					Subscribers: urls,
					Payload:     payload}); err != nil {
					g.lgs.WriteLog("ClientStream", err.Error())
				}
			} else {
				temp, err := connect(url)
				if err == nil {
					lock.Lock()
					stream, errConn = temp, err
					lock.Unlock()
				}

			}
		}
	}

	for _, url := range urls {
		cls := doQuery(url)
		arrayDoQuery = append(arrayDoQuery, cls)
	}

	for {
		select {
		case payload, ok := <-msg:
			if ok {
				for _, cls := range arrayDoQuery {
					go cls(payload)
				}
			}
		}
	}
}

func (g *grpcClient) ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, ctx context.Context, tx chan<- struct{}) {
	for _, url := range urls {
		go func(url string) {
			conn, _ := grpc.NewClient(url, g.opts...)
			client := pb.NewGrpcClient(conn)
			reqMsg := &pb.RequestMessage{
				Subscribers: nil,
				Payload:     initalMsg,
			}
			stream, err := client.ServerStream(ctx, reqMsg)
			if err == nil {
				for {
					resp, err := stream.Recv()
					if err != nil {
						g.lgs.WriteLog("ServerStream", err.Error())
						tx <- struct{}{}
						stream.CloseSend()
						return
					} else {
						msg <- resp.GetPayload()
					}
				}
			} else {
				return
			}

		}(url)
	}
}

func (g *grpcClient) BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- error, tx chan<- struct{}, ctx context.Context, done chan<- struct{}) {
	arrayMsg := make([]chan []byte, 0, len(urls))
	for p, url := range urls {
		ch := make(chan []byte, 0)
		arrayMsg = append(arrayMsg, ch)
		go func(ch <-chan []byte, p int, sub string) {
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
									c.CloseSend()
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

func NewGrpClient(opts []grpc.DialOption, lgs repos.LogsInt) repos.GrpClientInt {
	return &grpcClient{
		opts: opts,
		lgs:  lgs,
	}
}
