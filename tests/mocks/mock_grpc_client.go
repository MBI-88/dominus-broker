package mocks

import (
	"context"
	"dominus-project/app/domain/repos"
	"encoding/json"
	"time"
)

type grpcClientMock struct{}

func (grpcClientMock) Simple(url string, msg []byte) (repos.GrpResponseInt, error) {
	return nil, nil
}

func (grpcClientMock) ClientStream(urls []string, msg <-chan []byte, ctx context.Context) {
	for {
		select {
		case _, ok := <-msg:
			if !ok {
				return
			}
		}
	}

}

func (grpcClientMock) ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, ctx context.Context, tx chan<- struct{}) {
	for {
		select {
		case <-time.Tick(1 * time.Second):
			payload, _ := json.Marshal(Data)
			msg <- payload
		case <-time.Tick(5 * time.Second):
		case <-ctx.Done():
			for range urls {
				tx <- struct{}{}
			}
			return
		}
	}

}

func (grpcClientMock) BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- error, tx chan<- struct{}, ctx context.Context, done chan<- struct{}) {
	// received
	go func() {
		for {
			select {
			case _, ok := <-provMsg:
				if !ok {
					tx <- struct{}{}
					return
				}
			}
		}
	}()
	
	// subscribers sending messages
	for {
		select {
		case <-time.Tick(1 * time.Second):
			payload, _ := json.Marshal(Data)
			subMsg <- payload
		case <-time.Tick(5 * time.Second):
		case <-ctx.Done():
			for range urls {
				done <- struct{}{}
			}
			return
		}
	}

}

func NewGrpcClientMock() repos.GrpClientInt {
	return new(grpcClientMock)
}
