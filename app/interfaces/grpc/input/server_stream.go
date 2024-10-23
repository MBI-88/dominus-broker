package input


import (
	"context"
	"dominus/app/interactors"
	pb "dominus/app/interfaces/grpc/proto/builder"
)

type grpcContextServerStream struct {
	sr pb.Grpc_ServerStreamServer
}


func (g *grpcContextServerStream) Send(payload []byte) error {
	return g.sr.Send(&pb.ResponseMessage{
		Payload: payload,
	})
}

func (g *grpcContextServerStream) Context() context.Context {
	return g.sr.Context()
}


func newServerStreamContext(sr pb.Grpc_ServerStreamServer) interactors.StreamServerInt {
	return &grpcContextServerStream{
		sr: sr,
	}
}