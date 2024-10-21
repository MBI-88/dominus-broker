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
func (s *grpc) SendLocal(_ context.Context, ms *pb.Message) (*pb.Response, error) {
	var response *pb.Response

	return response, nil
}

//Receives array messages from client
func (s *grpc) SendClientStreamLocal(stream pb.API_SendClientStreamLocalClient) error {

	return nil
}

//Sends array messages to client
func (s *grpc) SendServerStreamLocal(ms *pb.Message, stream pb.API_SendServerStreamLocalServer) error {

	return nil
}

//Receives and sends messages from server to client
func (s *grpc) SendStreamLocal(stream pb.API_SendStreamLocalServer) error {

	return nil
}



func (s *grpc) SendRemote(_ context.Context, ms *pb.Message) (*pb.Response, error) {
	var response *pb.Response

	return response, nil
}

func (s *grpc) SendClientStreamRemote(stream pb.API_SendClientStreamRemoteClient) error {

	return nil
}


func (s *grpc) SendServerStreamRemote(ms *pb.Message, stream pb.API_SendServerStreamRemoteServer) error {

	return nil
}


func (s *grpc) SendStreamRemote(stream pb.API_SendStreamRemoteServer) error {

	return nil
}

