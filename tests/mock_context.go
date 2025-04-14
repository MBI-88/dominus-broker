package tests

import (
	"context"
	"dominus-project/app/domain/repos"
	"fmt"
	"time"
)


func helperErr(flag *simpleContextMock) {
	time.Sleep(3 * time.Second)
	ad.Lock()
	defer ad.Unlock()
	flag.err = fmt.Errorf("End")
}


type simpleContextMock struct {
	payload []byte
	subs []string
	err  error
}

func (s *simpleContextMock)  Descriptor() ([]byte, []int) {
	return nil, nil
}

func (s *simpleContextMock) GetPayload() []byte {
	return s.payload
}

func (s *simpleContextMock) GetSubscribers() []string {
	return s.subs
}

func (*simpleContextMock) Reset() {}

func (*simpleContextMock) String() string {
	return ""
}

func (*simpleContextMock) Validate() error {
	return nil
}

func (*simpleContextMock) ValidateAll() error {
	return nil
}

func NewSimpleContextMock(payload []byte, subs []string) repos.GrpRequestMessageInt {
	return &simpleContextMock {
		payload: payload,
		subs: subs,
		err: nil,
	}
}


type clientContextMock struct {
	*simpleContextMock
}

func (c *clientContextMock) Recv() (repos.GrpRequestMessageInt, error) {
	ad.Lock()
	defer ad.Unlock()
	return c.simpleContextMock, c.err
}


func NewClienContextMock(payload []byte, subs []string) repos.StreamClientInt {
	stream := &clientContextMock{
		&simpleContextMock{
			payload: payload,
			subs: subs,
			err: nil,
		},
	}
	go helperErr(stream.simpleContextMock)
	return stream
}


type serverContextMock struct {
	*simpleContextMock
}

func (se *serverContextMock) Send(payload []byte) error {
	ad.Lock()
	defer ad.Unlock()
	return se.err
}

func (se *serverContextMock) Context() context.Context {
	return context.Background()
}


func NewServerContextMock(payload []byte, subs []string) (repos.GrpRequestMessageInt,repos.StreamServerInt) {
	stream := &serverContextMock{
		&simpleContextMock{
			payload: payload,
			subs: subs,
			err: nil,
		},
	}
	go helperErr(stream.simpleContextMock)
	return stream.simpleContextMock, stream
}



type biContextMock struct {
	*simpleContextMock
}

func (b *biContextMock) Recv() (repos.GrpRequestMessageInt, error) {
	ad.Lock()
	defer ad.Unlock()
	return b.simpleContextMock, b.err
}

func (b *biContextMock) Send(msg []byte) error {
	ad.Lock()
	defer ad.Unlock()
	return b.err
}


func NewBiContextMock(payload []byte, subs []string) repos.StreamBiInt {
	stream := &biContextMock{
		&simpleContextMock{
			payload: payload,
			subs: subs,
			err: nil,
		},
	}
	go helperErr(stream.simpleContextMock)
	return stream
}