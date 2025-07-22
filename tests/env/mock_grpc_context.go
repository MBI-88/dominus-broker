package env

import (
	"context"
	"dominus-project/internal/domain/adapters"
	"fmt"
	"time"
)

func helperErr(flag *simpleContextMock) {
	time.Sleep(3 * time.Second)
	Ad.Lock()
	defer Ad.Unlock()
	flag.err = fmt.Errorf("End")
}

type simpleContextMock struct {
	payload []byte
	subs    []string
	err     error
}

func NewSimpleContextMock(payload []byte, subs []string) adapters.IGrpcDto {
	return &simpleContextMock{
		payload: payload,
		subs:    subs,
		err:     nil,
	}
}

func (s *simpleContextMock) GetPayload() []byte {
	return s.payload
}

func (s *simpleContextMock) GetSubscribers() []string {
	return s.subs
}

type clientContextMock struct {
	*simpleContextMock
}

func NewClienContextMock(payload []byte, subs []string) adapters.IStreamClient {
	stream := &clientContextMock{
		&simpleContextMock{
			payload: payload,
			subs:    subs,
			err:     nil,
		},
	}
	go helperErr(stream.simpleContextMock)
	return stream
}

func (c *clientContextMock) Recv() (adapters.IGrpcDto, error) {
	Ad.Lock()
	defer Ad.Unlock()
	return c.simpleContextMock, c.err
}

type serverContextMock struct {
	*simpleContextMock
}

func NewServerContextMock(payload []byte, subs []string) (adapters.IGrpcDto, adapters.IStreamServer) {
	stream := &serverContextMock{
		&simpleContextMock{
			payload: payload,
			subs:    subs,
			err:     nil,
		},
	}
	go helperErr(stream.simpleContextMock)
	return stream.simpleContextMock, stream
}

func (se *serverContextMock) Send(payload []byte) error {
	Ad.Lock()
	defer Ad.Unlock()
	return se.err
}

func (se *serverContextMock) Context() context.Context {
	return context.Background()
}

type biContextMock struct {
	*simpleContextMock
}

func NewBiContextMock(payload []byte, subs []string) adapters.IStreamBi {
	stream := &biContextMock{
		&simpleContextMock{
			payload: payload,
			subs:    subs,
			err:     nil,
		},
	}
	go helperErr(stream.simpleContextMock)
	return stream
}

func (b *biContextMock) Recv() (adapters.IGrpcDto, error) {
	Ad.Lock()
	defer Ad.Unlock()
	return b.simpleContextMock, b.err
}

func (b *biContextMock) Send(msg []byte) error {
	Ad.Lock()
	defer Ad.Unlock()
	return b.err
}