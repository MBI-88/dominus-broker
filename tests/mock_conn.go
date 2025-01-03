package tests

import (
	"context"
	"dominus/app/interactors"
	"fmt"
	"time"
)


func helperErr(flag *simpleconnMock) {
	time.Sleep(5 * time.Second)
	adquired.Lock()
	defer adquired.Unlock()
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

func NewSimpleConnMock(payload []byte, subs []string) interactors.GrpRequestMessageInt {
	return &simpleconnMock {
		payload: payload,
		subs: subs,
		err: nil,
	}
}


type clientconnMock struct {
	*simpleconnMock
}

func (c *clientconnMock) Recv() (interactors.GrpRequestMessageInt, error) {
	return c.simpleconnMock, c.err
}


func NewClienConnMock(payload []byte, subs []string) interactors.StreamClientInt {
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
	return se.err
}

func (se *serverconnMock) Context() context.Context {
	return context.Background()
}


func NewServerConnMock(payload []byte, subs []string) (interactors.GrpRequestMessageInt,interactors.StreamServerInt) {
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

func (b *biconnMock) Recv() (interactors.GrpRequestMessageInt, error) {
	return b.simpleconnMock, b.err
}

func (b *biconnMock) Send(msg []byte) error {
	return b.err
}


func NewBiConnMock(payload []byte, subs []string) interactors.StreamBiInt {
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