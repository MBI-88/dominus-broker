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

func (s *simpleContextMock) Descriptor() ([]byte, []int) {
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

func NewSimpleContextMock(payload []byte, subs []string) adapters.IGrpcDto {
	return &simpleContextMock{
		payload: payload,
		subs:    subs,
		err:     nil,
	}
}

type clientContextMock struct {
	*simpleContextMock
}

func (c *clientContextMock) Recv() (adapters.IGrpcDto, error) {
	Ad.Lock()
	defer Ad.Unlock()
	return c.simpleContextMock, c.err
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

type serverContextMock struct {
	*simpleContextMock
}

func (se *serverContextMock) Send(payload []byte) error {
	Ad.Lock()
	defer Ad.Unlock()
	return se.err
}

func (se *serverContextMock) Context() context.Context {
	return context.Background()
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

type biContextMock struct {
	*simpleContextMock
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
