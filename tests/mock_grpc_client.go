package tests

import (
	"dominus/app/domain/entities"
	"dominus/app/interactors"
	"encoding/json"
	"sync"
	"time"
)

type grpcClientMock struct {
}

func (*grpcClientMock) Simple(url string, msg []byte) (interactors.GrpResponseInt, error) {
	return nil, nil
}

func (*grpcClientMock) ClientStream(urls []string, msg <-chan []byte, sig chan<- entities.Logs, tx chan<- struct{}) {
	for {
		select {
		case payload, ok := <-msg:
			if ok {
				println(string(payload))
			} else {
				tx <- struct{}{}
				return
			}
		}
	}

}

func (*grpcClientMock) ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, sig chan<- entities.Logs, closed <-chan struct{}, tx chan<- struct{}) {
	open := new(bool)
	*open = true
	sync := new(sync.Mutex)
	go func() {
		<-closed
		sync.Lock()
		defer sync.Unlock()
		*open = false
	}()
	for {
		select {
		case <-time.Tick(1 * time.Second):
			payload, _ := json.Marshal(data)
			if *open {
				msg <- payload
			} else {
				for range urls {
					tx <- struct{}{}
				}
			}
		case <-time.Tick(10 * time.Second):
			for range urls {
				tx <- struct{}{}
			}
			return
		}
	}

}

func (*grpcClientMock) BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- entities.Logs, tx chan<- struct{}, rx <-chan struct{}, done chan<- struct{}) {

}

func NewGrpcClientMock() interactors.GrpClientInt {
	return new(grpcClientMock)
}
