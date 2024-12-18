package input


import (
	"dominus/app/interactors"
	pb "dominus/app/interfaces/grpc/proto/builder"
)


type grpcBiContextStream struct {
	sr pb.Grpc_BidirectionalStreamServer
}


func (g *grpcBiContextStream) Recv() (interactors.GrpRequestMessageInt, error) {
	return g.sr.Recv()
}

func (g *grpcBiContextStream) Send(msg []byte) error {
	return g.sr.Send(&pb.ResponseMessage{
		Payload: msg,
	})
}


func newBiStreamConn(sr pb.Grpc_BidirectionalStreamServer) interactors.StreamBiInt {
	return &grpcBiContextStream{
		sr: sr,
	}
}