package input

import (
	"context"
	"dominus/app/interactors"
	pb "dominus/app/interfaces/grpc/proto/builder"

	jsoniter "github.com/json-iterator/go"
)


type grpc struct {
	pb.UnimplementedAPIServer
	inter interactors.InteractorInt
	js    jsoniter.API
}

//Receives simple messages from client
func (s *grpc) Send(_ context.Context, ms *pb.RequestMessage) (*pb.Response, error) {
	var response *pb.Response

	return response, nil
}

//Receives array messages from client
func (s *grpc) SendClientStream(stream pb.API_SendServerStreamClient) error {

	return nil
}

//Sends array messages to client
func (s *grpc) RxClientStream(stream pb.API_SendClientStreamServer) error {

	return nil
}

//Receives and sends messages from server to client
func (s *grpc) SendfullStream(stream pb.API_SendFullDuplexStreamServer) error {

	return nil
}
