package tests

import (
	"context"
	"dominus/app/domain/repos"
	"fmt"
	"time"
)


func helperErr(flag *simpleconnMock) {
	time.Sleep(3 * time.Second)
	ad.Lock()
	defer ad.Unlock()
	flag.err = fmt.Errorf("End")
}


type simpleconnMock struct {
	payload []byte
	subs []string
	err  error
}

func (s *simpleconnMock)  Descriptor() ([]byte, []int) {
	return nil, nil
}

func (s *simpleconnMock) GetPayload() []byte {
	return s.payload
}

func (s *simpleconnMock) GetSubscribers() []string {
	return s.subs
}

func (*simpleconnMock) Reset() {}

func (*simpleconnMock) String() string {
	return ""
}

func (*simpleconnMock) Validate() error {
	return nil
}

func (*simpleconnMock) ValidateAll() error {
	return nil
}

func NewSimpleConnMock(payload []byte, subs []string) repos.GrpRequestMessageInt {
	return &simpleconnMock {
		payload: payload,
		subs: subs,
		err: nil,
	}
}


type clientconnMock struct {
	*simpleconnMock
}

func (c *clientconnMock) Recv() (repos.GrpRequestMessageInt, error) {
	ad.Lock()
	defer ad.Unlock()
	return c.simpleconnMock, c.err
}


func NewClienConnMock(payload []byte, subs []string) repos.StreamClientInt {
	stream := &clientconnMock{
		&simpleconnMock{
			payload: payload,
			subs: subs,
			err: nil,
		},
	}
	go helperErr(stream.simpleconnMock)
	return stream
}


type serverconnMock struct {
	*simpleconnMock
}

func (se *simpleconnMock) Send(payload []byte) error {
	ad.Lock()
	defer ad.Unlock()
	return se.err
}

func (se *serverconnMock) Context() context.Context {
	return context.Background()
}


func NewServerConnMock(payload []byte, subs []string) (repos.GrpRequestMessageInt,repos.StreamServerInt) {
	stream := &serverconnMock{
		&simpleconnMock{
			payload: payload,
			subs: subs,
			err: nil,
		},
	}
	go helperErr(stream.simpleconnMock)
	return stream.simpleconnMock, stream
}



type biconnMock struct {
	*simpleconnMock
}

func (b *biconnMock) Recv() (repos.GrpRequestMessageInt, error) {
	ad.Lock()
	defer ad.Unlock()
	return b.simpleconnMock, b.err
}

func (b *biconnMock) Send(msg []byte) error {
	ad.Lock()
	defer ad.Unlock()
	return b.err
}


func NewBiConnMock(payload []byte, subs []string) repos.StreamBiInt {
	stream := &biconnMock{
		&simpleconnMock{
			payload: payload,
			subs: subs,
			err: nil,
		},
	}
	go helperErr(stream.simpleconnMock)
	return stream
}