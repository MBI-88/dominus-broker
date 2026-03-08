package output

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/domain/enum"
	"io"
	"sync"
	"time"

	pb "github.com/PR0C0D3-MBI/dominus-proto-definition/dominus"

	"google.golang.org/grpc"
)

type brokerClient struct {
	opts []grpc.DialOption
	lgs  adapters.Logs
}

func NewGrpClient(opts []grpc.DialOption, lgs adapters.Logs) adapters.BrokerClient {
	return &brokerClient{
		opts: opts,
		lgs:  lgs,
	}
}

func (g *brokerClient) ClientStream(urls []string, msg <-chan []byte, ctx context.Context) {
	arrayDoQuery := make([]func([]byte), 0, len(urls))
	lock := new(sync.Mutex)
	connect := func(url string) (pb.API_ClientStreamClient, error) {
		conn, _ := grpc.NewClient(url, g.opts...)
		c := pb.NewAPIClient(conn)
		return c.ClientStream(ctx)
	}
	doQuery := func(url string) func([]byte) {
		stream, errConn := connect(url)
		return func(payload []byte) {
			go g.lgs.WriteLog(ctx, enum.DEBUG, "ClientStream", "debuging sending request")
			lock.Lock()
			err := errConn
			lock.Unlock()
			if err == nil {
				lock.Lock()
				err := stream.Send(&pb.StreamRequestMessage{
					Subscribers: urls,
					Payload:     payload})
				lock.Unlock()
				if err != nil {
					go g.lgs.WriteLog(ctx , enum.ERROR ,"ClientStream", err.Error())
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

	for m := range msg {
		for _, cls := range arrayDoQuery {
			go cls(m)
		}
	}
}

func (g *brokerClient) ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, ctx context.Context, tx chan<- struct{}) {
	for _, url := range urls {
		go func(url string) {
			conn, _ := grpc.NewClient(url, g.opts...)
			client := pb.NewAPIClient(conn)
			reqMsg := &pb.StreamRequestMessage{
				Subscribers: nil,
				Payload:     initalMsg,
			}
		connect:
			stream, err := client.ServerStream(ctx, reqMsg)
			if err == nil {
				for {
					resp, err := stream.Recv()
					if err == io.EOF {
						go g.lgs.WriteLog(ctx, enum.DEBUG ,"ServerStream", err.Error())
						tx <- struct{}{}
						if err := stream.CloseSend(); err != nil {
							go g.lgs.WriteLog(ctx, enum.DEBUG ,"ColseSend", err.Error())
						}
						return
					} else if err != io.EOF && err != nil {
						go g.lgs.WriteLog(ctx, enum.DEBUG, "ServerStream EOF", err.Error())
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

func (g *brokerClient) BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- error, tx chan<- struct{}, ctx context.Context, done chan<- struct{}) {
	lock := new(sync.Mutex)
	arrayDoQuery := make([]func([]byte), 0, len(urls))

	connect := func(url string) (pb.API_BidirectionalStreamClient, error) {
		conn, _ := grpc.NewClient(url, g.opts...)
		c := pb.NewAPIClient(conn)
		return c.BidirectionalStream(ctx)
	}

	doQuey := func(url string, stream pb.API_BidirectionalStreamClient, errConn error) func([]byte) {
		return func(payload []byte) {
			go g.lgs.WriteLog(ctx, enum.DEBUG, "BidirectionalStream", "debuging sending request")
			lock.Lock()
			err := errConn
			lock.Unlock()
			if err == nil {
				lock.Lock()
				err := stream.Send(&pb.StreamRequestMessage{
					Subscribers: urls,
					Payload:     payload})
				lock.Unlock()
				if err != nil {
					go g.lgs.WriteLog(ctx,enum.ERROR ,"BidirectionalStream", err.Error())
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

		go func(url string, c pb.API_BidirectionalStreamClient, errConn error) {
		connect:
			go g.lgs.WriteLog(ctx, enum.DEBUG, "Connect", "debuging sending request")
			lock.Lock()
			err := errConn
			lock.Unlock()
			if err == nil {
				for {
					lock.Lock()
					resp, err := stream.Recv()
					lock.Unlock()
					if err == io.EOF {
						go g.lgs.WriteLog(ctx, enum.ERROR ,"ServerStream", err.Error())
						tx <- struct{}{}
						if err := stream.CloseSend(); err != nil {
							go g.lgs.WriteLog(ctx, enum.ERROR, "CloseSend", err.Error())
						}
						return
					} else if err != io.EOF && err != nil {
						go g.lgs.WriteLog(ctx, enum.DEBUG, "Connect", "reconecting ...")
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

	for msg := range provMsg {
		for _, cls := range arrayDoQuery {
			go cls(msg)
		}
	}
	tx <- struct{}{}
}
