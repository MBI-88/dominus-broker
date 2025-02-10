package tests

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/repos"
	"encoding/json"
	"sync"
	"time"
)

type grpcClientMock struct{}

func (grpcClientMock) Simple(url string, msg []byte) (repos.GrpResponseInt, error) {
	return nil, nil
}

func (grpcClientMock) ClientStream(urls []string, msg <-chan []byte, sig chan<- entities.Logs, tx chan<- struct{}) {
	for {
		select {
		case _, ok := <-msg:
			if !ok {
				tx <- struct{}{}
				return
			}
		}
	}

}

func (grpcClientMock) ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, sig chan<- entities.Logs, closed <-chan struct{}, tx chan<- struct{}) {
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
				return
			}
		case <-time.Tick(10 * time.Second):
			for range urls {
				tx <- struct{}{}
			}
			return
		}
	}

}

func (grpcClientMock) BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- entities.Logs, tx chan<- struct{}, rx <-chan struct{}, done chan<- struct{}) {
	open := new(bool)
	*open = true
	sync := new(sync.Mutex)
	go func() {
		<-rx
		sync.Lock()
		defer sync.Unlock()
		*open = false
	}()

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

	for {
		select {
		case <-time.Tick(1 * time.Second):
			payload, _ := json.Marshal(data)
			if *open {
				subMsg <- payload
			} else {
				for range urls {
					done <- struct{}{}
				}
				return
			}
		case <-time.Tick(10 * time.Second):
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
