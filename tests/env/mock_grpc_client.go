package env

import (
	"context"
	"dominus-project/internal/domain/repos"
	"encoding/json"
	"fmt"
	"time"
)

type mockResponse struct {
	happyPath bool
}

func (mockResponse) Descriptor() ([]byte, []int) {
	return  []byte("empty"), []int{}
}
func (mockResponse) GetMessage() string {
	return  ""
}
func (m *mockResponse) GetStatus() uint32  {
	if m.happyPath {
		return  200
	}
	return  400
}
func (mockResponse) ProtoMessage() {}
func (mockResponse) Reset() {}
func (mockResponse) String() string {
	return  ""
}
func (mockResponse) Validate() error {
	return  nil
}
func (mockResponse) ValidateAll() error {
	return  nil
}



type grpcClientMock struct{
	rspo repos.IGrpResponse
	happyPath bool
}

func (g *grpcClientMock) Simple(url string, msg []byte) (repos.IGrpResponse, error) {
	if g.happyPath {
		return g.rspo, nil
	}
	return  nil, fmt.Errorf("Error")
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

func NewGrpcClientMock(happyPathCls, happyPathResp bool) repos.IGrpClient {
	return &grpcClientMock{
		rspo: &mockResponse{
			happyPath: happyPathResp,
		},
		happyPath: happyPathCls,
	}
}
