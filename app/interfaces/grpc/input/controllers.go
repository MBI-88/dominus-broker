package input

import (
	"context"
	"dominus/app/interactors"
	pb "dominus/app/interfaces/grpc/proto/builder"

	"google.golang.org/grpc"
)


type grpcController struct {
	pb.UnimplementedGrpcServer
	inter interactors.InteractorInt
}

//Receives simple messages from client
func (s *grpcController) Simple(_ context.Context, ms *pb.RequestMessage) (*pb.Response, error) {
	inter := s.inter.NewConnection()

	if err  := inter.SimpleConn(ms); err != nil {
		return &pb.Response{Status: uint32(406),Message: err.Error()}, err
	}

	return &pb.Response{Status: uint32(406),Message: "Accepted"}, nil
}

//Receives array messages from client
func (s *grpcController) ClientStream(stream pb.Grpc_ClientStreamServer) error {
	inter := s.inter.NewConnection()
	ctx := newClientStreamContext(stream)

	if err := inter.StreamClientConn(ctx); err != nil {
		stream.SendAndClose(&pb.Response{
			Status: uint32(500),
			Message: err.Error(),
		})
		return err
	}
	return nil
}

//Sends array messages to client
func (s *grpcController) ServerStream(ms *pb.RequestMessage, stream pb.Grpc_ServerStreamServer) error {
	inter := s.inter.NewConnection()
	ctx := newServerStreamContext(stream)
	if err := inter.StreamServerConn(ms,ctx); err != nil {
		return err
	}
	return nil
}

//Receives and sends messages from server to client
func (s *grpcController) BidirectionalStream(stream pb.Grpc_BidirectionalStreamServer) error {
	inter := s.inter.NewConnection()
	ctx := newBiStreamConn(stream)
	if err := inter.StreamBiConn(ctx); err != nil {
		return err
	}
	return nil
}


func NewGrpcServe(opts []grpc.ServerOption, i interactors.InteractorInt) *grpc.Server {
	s := grpc.NewServer(opts...)
	pb.RegisterGrpcServer(s, &grpcController{inter:i})
	return s
}