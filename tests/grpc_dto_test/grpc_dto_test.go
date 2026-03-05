package grpcdto_test

import (
	"context"
	"dominus-project/internal/infrastructure/grpc/dto"
	"testing"

	"github.com/PR0C0D3-MBI/dominus-proto-definition/dominus"
	"google.golang.org/grpc/metadata"
)

type receiveGrpc struct {
	payload []byte
	subs    []string
	ctx     context.Context
}

func (r *receiveGrpc) Recv() (*dominus.StreamRequestMessage, error) {
	return &dominus.StreamRequestMessage{
		Subscribers: r.subs,
		Payload:     r.payload,
	}, nil
}

func (r *receiveGrpc) Context() context.Context {
	return r.ctx
}

func (r *receiveGrpc) RecvMsg(d any) error {

	return nil
}

func (r *receiveGrpc) SendAndClose(*dominus.StreamResponseMessage) error {
	return nil
}

func (r *receiveGrpc) SendHeader(metadata.MD) error {

	return nil
}

func (r *receiveGrpc) SendMsg(d any) error {
	return nil
}

func (r *receiveGrpc) SetHeader(metadata.MD) error {

	return nil
}

func (r *receiveGrpc) SetTrailer(metadata.MD) {

}

func (r *receiveGrpc) Send(d *dominus.StreamResponseMessage) error {

	return nil
}

func TestGrpcDto(t *testing.T) {

	t.Run("ClientContext ok", func(t *testing.T) {
		p := []byte("test")
		sr := &receiveGrpc{
			payload: p,
			subs:    []string{"server1.api.com", "server2.api.com"},
		}
		dto := dto.NewClientStreamContext(sr)

		data, err := dto.Recv()
		if err != nil {
			t.Fatal(err)
		}
		resp := data.GetPayload()
		if string(resp) != string(p) {
			t.Fatalf("Expected %s received %s", resp, "test")
		}

		subs := data.GetSubscribers()
		if len(subs) == 0 {
			t.Fatalf("Expecte subscribers > 0")
		}
	})

	t.Run("ServerStreamContext ok", func(t *testing.T) {
		p := []byte("test")
		sr := &receiveGrpc{
			payload: p,
			subs:    []string{"server1.api.com", "server2.api.com"},
		}
		dto := dto.NewServerStreamContext(sr)

		err := dto.Send(p)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("BidirectionalStreamContext", func(t *testing.T) {
		p := []byte("test")
		sr := &receiveGrpc{
			payload: p,
			subs:    []string{"server1.api.com", "server2.api.com"},
		}
		dto := dto.NewBiStreamConn(sr)

		err := dto.Send(p)
		if err != nil {
			t.Fatal(err)
		}

		data, err := dto.Recv()
		if err != nil {
			t.Fatal(err)
		}
		resp := data.GetPayload()
		if string(resp) != string(p) {
			t.Fatalf("Expected %s received %s", resp, "test")
		}

		subs := data.GetSubscribers()
		if len(subs) == 0 {
			t.Fatalf("Expecte subscribers > 0")
		}
	})

	t.Run("Context ok", func(t *testing.T) {
		p := []byte("test")
		sr := &receiveGrpc{
			payload: p,
			subs:    []string{"server1.api.com", "server2.api.com"},
			ctx:     context.Background(),
		}
		dto := dto.NewServerStreamContext(sr)

		ctx := dto.Context()

		if ctx == nil {
			t.Fatalf("Expected output different from nil")
		}

	})

}
