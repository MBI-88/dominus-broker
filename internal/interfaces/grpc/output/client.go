package output

import (
	"context"
	"dominus-project/internal/domain/repos"
	pb "dominus-project/internal/interfaces/grpc/proto/builder"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc"
)

type grpcClient struct {
	opts []grpc.DialOption
	lgs  repos.ILogs
}

func (g *grpcClient) Simple(url string, body []byte) (repos.IGrpResponse, error) {
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

func (g *grpcClient) ClientStream(urls []string, msg <-chan []byte, ctx context.Context) {
	arrayDoQuery := make([]func([]byte), 0, len(urls))
	lock := new(sync.Mutex)
	connect := func(url string) (pb.Grpc_ClientStreamClient, error) {
		conn, _ := grpc.NewClient(url, g.opts...)
		c := pb.NewGrpcClient(conn)
		return c.ClientStream(ctx)
	}
	doQuery := func(url string) func([]byte) {
		stream, errConn := connect(url)
		return func(payload []byte) {
			lock.Lock()
			err := errConn
			lock.Unlock()
			if err == nil {
				lock.Lock()
				err := stream.Send(&pb.RequestMessage{
					Subscribers: urls,
					Payload:     payload})
				lock.Unlock()
				if err != nil {
					g.lgs.WriteLog("ClientStream", err.Error())
					lock.Lock()
					errConn = err
					lock.Unlock()
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
		connect:
			stream, err := client.ServerStream(ctx, reqMsg)
			if err == nil {
				for {
					resp, err := stream.Recv()
					if err == io.EOF {
						g.lgs.WriteLog("ServerStream", err.Error())
						tx <- struct{}{}
						stream.CloseSend()
						return
					} else if err != io.EOF && err != nil {
						goto connect
					} else {
						msg <- resp.GetPayload()
					}
				}
			} else {
				for {
					select {
					case <-ctx.Done():
						return
					case <-time.Tick(1 * time.Millisecond):
						goto connect
					}
				}
			}
		}(url)
	}
}

func (g *grpcClient) BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- error, tx chan<- struct{}, ctx context.Context, done chan<- struct{}) {
	lock := new(sync.Mutex)
	arrayDoQuery := make([]func([]byte), 0, len(urls))

	connect := func(url string) (pb.Grpc_BidirectionalStreamClient, error) {
		conn, _ := grpc.NewClient(url, g.opts...)
		c := pb.NewGrpcClient(conn)
		return c.BidirectionalStream(ctx)
	}

	doQuey := func(url string, stream pb.Grpc_BidirectionalStreamClient, errConn error) func([]byte) {
		return func(payload []byte) {
			lock.Lock()
			err := errConn
			lock.Unlock()
			if err == nil {
				lock.Lock()
				err := stream.Send(&pb.RequestMessage{
					Subscribers: urls,
					Payload:     payload})
				lock.Unlock()
				if err != nil {
					g.lgs.WriteLog("ClientStream", err.Error())
					lock.Lock()
					errConn = err
					lock.Unlock()
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
		stream, err := connect(url)
		cls := doQuey(url, stream, err)
		arrayDoQuery = append(arrayDoQuery, cls)

		go func(url string, c pb.Grpc_BidirectionalStreamClient, errConn error) {
		connect:
			lock.Lock()
			err := errConn
			lock.Unlock()
			if err == nil {
				for {
					lock.Lock()
					resp, err := stream.Recv()
					lock.Unlock()
					if err == io.EOF {
						g.lgs.WriteLog("ServerStream", err.Error())
						tx <- struct{}{}
						stream.CloseSend()
						return
					} else if err != io.EOF && err != nil {
						goto connect
					} else {
						subMsg <- resp.GetPayload()
					}
				}
			} else {
				for {
					select {
					case <-ctx.Done():
						return
					case <-time.Tick(1 * time.Millisecond):
						goto connect
					}
				}
			}
		}(url, stream, err)
	}

	for {
		select {
		case payload, ok := <-provMsg:
			if ok {
				for _, cls := range arrayDoQuery {
					go cls(payload)
				}
			} else {
				tx <- struct{}{}
				return
			}
		}
	}
}

func NewGrpClient(opts []grpc.DialOption, lgs repos.ILogs) repos.IGrpClient {
	return &grpcClient{
		opts: opts,
		lgs:  lgs,
	}
}
