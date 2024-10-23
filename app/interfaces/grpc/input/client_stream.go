package input

import (
	"dominus/app/interactors"
	pb "dominus/app/interfaces/grpc/proto/builder"
)

type grpcContextClientStream struct {
	sr pb.Grpc_ClientStreamServer
}

func (g *grpcContextClientStream) Recv() (interactors.GrpRequestMessageInt, error) {
	return g.sr.Recv()
}



func newClientStreamContext(sr pb.Grpc_ClientStreamServer) interactors.StreamClientInt {
	return &grpcContextClientStream{
		sr: sr,
	}
}
