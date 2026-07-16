package outbound

import (
	"context"
	"dominus-broker/internal/domain/repositories"
	"dominus-broker/internal/infrastructure/enum"
	"dominus-broker/internal/infrastructure/event"
	"io"
	"sync"
	"time"

	pb "github.com/MBI-88/dominus-proto-definition/dominus"

	"google.golang.org/grpc"
)

type clientStream struct {
	opts []grpc.DialOption
	lgs  event.Event
}

func NewGrpClient(opts []grpc.DialOption, lgs event.Event) repositories.BrokerClient {
	return &clientStream{
		opts: opts,
		lgs:  lgs,
	}
}

func (g *clientStream) ClientStream(urls []string, msg <-chan []byte, ctx context.Context) {
	arrayDoQuery := make([]func([]byte), 0, len(urls))
	lock := new(sync.Mutex)

	connect := func(url string) (pb.BrokerAPI_ClientStreamClient, error) {
		conn, err := grpc.NewClient(url, g.opts...)
		if err != nil {
			return nil, err
		}
		c := pb.NewBrokerAPIClient(conn)
		return c.ClientStream(ctx)
	}
	doQuery := func(url string) func([]byte) {
		stream, errConn := connect(url)
		return func(payload []byte) {
			g.lgs.WriteLog(ctx, enum.DEBUG, "outbound.ClientStream.doQuery", enum.DEBUG_DESCRIPTION)
			lock.Lock()
			err := errConn
			lock.Unlock()
			if err == nil {
				lock.Lock()
				err := stream.Send(&pb.StreamRequestMessage{
					Subscribers: nil,
					Payload:     payload})
				lock.Unlock()
				if err != nil {
					g.lgs.WriteLog(ctx, enum.ERROR, "outbound.ClientStream.Send", err.Error())
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

func (g *clientStream) ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, ctx context.Context, tx chan<- struct{}) {
	var workerWG sync.WaitGroup
	workerWG.Add(len(urls))

	connect := func(url string) (pb.BrokerAPI_ServerStreamClient, error) {
		conn, err := grpc.NewClient(url, g.opts...)
		if err != nil {
			return nil, err
		}
		client := pb.NewBrokerAPIClient(conn)
		reqMsg := &pb.StreamRequestMessage{
			Subscribers: nil,
			Payload:     initalMsg,
		}
		return client.ServerStream(ctx, reqMsg)
	}

	for _, url := range urls {
		go func(url string) {
			defer workerWG.Done()
		connectLabel:
			stream, err := connect(url)
			if err == nil {
				for {
					resp, err := stream.Recv()
					if err == io.EOF {
						g.lgs.WriteLog(ctx, enum.DEBUG, "outbound.ServerStream.Recv", enum.DEBUG_DESCRIPTION)
						return
					} else if err != io.EOF && err != nil {
						g.lgs.WriteLog(ctx, enum.ERROR, "outbound.ServerStream.EOF", err.Error())
						goto connectLabel
					} else {
						msg <- resp.GetPayload()
					}
				}
			} else {
				retry := time.NewTicker(5 * time.Millisecond) // retrying
				for {
					select {
					case <-ctx.Done():
						retry.Stop()
						return
					case <-retry.C:
						retry.Stop()
						goto connectLabel
					}
				}
			}
		}(url)
	}

	workerWG.Wait()
	if tx != nil {
		tx <- struct{}{}
	}

}

func (g *clientStream) BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, tx chan<- struct{}, ctx context.Context) {
	lock := new(sync.Mutex)
	type endpoint struct {
		url    string
		stream pb.BrokerAPI_BidirectionalStreamClient
		err    error
	}

	connect := func(url string) (pb.BrokerAPI_BidirectionalStreamClient, error) {
		conn, err := grpc.NewClient(url, g.opts...)
		if err != nil {
			return nil, err
		}
		c := pb.NewBrokerAPIClient(conn)
		return c.BidirectionalStream(ctx)
	}

	total := len(urls)
	arrayDoQuery := make([]func([]byte), 0, total)
	var workerWG sync.WaitGroup
	workerWG.Add(total)

	for _, u := range urls {
		stream, err := connect(u)
		ep := &endpoint{url: u, stream: stream, err: err}

		cls := func(payload []byte) {
			g.lgs.WriteLog(ctx, enum.DEBUG, "outbound.BidirectionalStream.doQuery", enum.DEBUG_DESCRIPTION)
			lock.Lock()
			st := ep.stream
			e := ep.err
			lock.Unlock()
			if e == nil && st != nil {
				serr := st.Send(&pb.StreamRequestMessage{
					Subscribers: nil,
					Payload:     payload,
				})
				if serr != nil {
					g.lgs.WriteLog(ctx, enum.ERROR, "outbound.BidirectionalStream.doQuery", serr.Error())
					lock.Lock()
					ep.err = serr
					lock.Unlock()
				}
			} else {
				temp, cerr := connect(u)
				if cerr == nil {
					lock.Lock()
					ep.stream, ep.err = temp, nil
					lock.Unlock()
				}
			}
		}
		arrayDoQuery = append(arrayDoQuery, cls)

		go func(ep *endpoint) {
			defer workerWG.Done()

		connectLabel:
			g.lgs.WriteLog(ctx, enum.DEBUG, "outbound.BidirectionalStream.Connect", enum.DEBUG_DESCRIPTION)
			lock.Lock()
			e := ep.err
			st := ep.stream
			lock.Unlock()

			if e != nil || st == nil {
				retry := time.NewTicker(5 * time.Millisecond)
				for {
					select {
					case <-ctx.Done():
						retry.Stop()
						return
					case <-retry.C:
						retry.Stop()
						temp, cerr := connect(ep.url)
						if cerr == nil {
							lock.Lock()
							ep.stream, ep.err = temp, nil
							lock.Unlock()
						}
						goto connectLabel
					}
				}
			}

			for {
				resp, rerr := st.Recv()
				if rerr == io.EOF {
					g.lgs.WriteLog(ctx, enum.DEBUG, "outbound.BidirectionalStream.Recv", enum.DEBUG_DESCRIPTION)
					return
				}
				if rerr != nil && rerr != io.EOF {
					if ctx.Err() != nil {
						g.lgs.WriteLog(ctx, enum.ERROR, "outbound.BidirectionalStream.ctx", ctx.Err().Error())
						return
					}
					lock.Lock()
					ep.err = rerr
					lock.Unlock()
					goto connectLabel
				}
				select {
				case <-ctx.Done():
					return
				case subMsg <- resp.GetPayload():
				}
			}
		}(ep)
	}

	for msg := range provMsg {
		for _, cls := range arrayDoQuery {
			go cls(msg)
		}
	}

	workerWG.Wait()
	if tx != nil {
		tx <- struct{}{}
	}
}
