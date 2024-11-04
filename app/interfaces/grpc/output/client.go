package output

import (
	"context"
	"dominus/app/interactors"
	pb "dominus/app/interfaces/grpc/proto/builder"

	"google.golang.org/grpc"
)


type grpcClient struct {
	opts []grpc.DialOption

}


func (g *grpcClient) Simple(url string, body []byte ) (interactors.GrpResponseInt, error) {
	conn, err := grpc.NewClient(url, g.opts...)
	if err != nil  {
		return nil, err
	}

	client := pb.NewGrpcClient(conn)
	ctx := context.Background()

	msg := &pb.RequestMessage{
		Topic: "", Payload: body,
	}

	resp, err := client.Simple(ctx, msg)
	if err != nil {
		return nil,  err
	}
	
	return resp, nil 
}

func (g *grpcClient) ClientStream(url string) {

}

func (g *grpcClient) ServerStream(url string) {

}

func (g *grpcClient) BidirectionalStream(url string) {

}




func NewGrpClient(opts []grpc.DialOption) interactors.GrpClientInt {
	

	return &grpcClient{}
}