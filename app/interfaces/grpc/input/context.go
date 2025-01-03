package input

import (
	"context"
	"dominus/app/domain/repos"
	pb "dominus/app/interfaces/grpc/proto/builder"
)


type grpcBiContextStream struct {
	sr pb.Grpc_BidirectionalStreamServer
}


func (g *grpcBiContextStream) Recv() (repos.GrpRequestMessageInt, error) {
	return g.sr.Recv()
}

func (g *grpcBiContextStream) Send(msg []byte) error {
	return g.sr.Send(&pb.ResponseMessage{
		Payload: msg,
	})
}


func newBiStreamConn(sr pb.Grpc_BidirectionalStreamServer) repos.StreamBiInt {
	return &grpcBiContextStream{
		sr: sr,
	}
}




type grpcContextClientStream struct {
	sr pb.Grpc_ClientStreamServer
}

func (g *grpcContextClientStream) Recv() (repos.GrpRequestMessageInt, error) {
	return g.sr.Recv()
}



func newClientStreamContext(sr pb.Grpc_ClientStreamServer) repos.StreamClientInt {
	return &grpcContextClientStream{
		sr: sr,
	}
}





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


func newServerStreamContext(sr pb.Grpc_ServerStreamServer) repos.StreamServerInt {
	return &grpcContextServerStream{
		sr: sr,
	}
}