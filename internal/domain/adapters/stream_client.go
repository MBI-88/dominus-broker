package adapters

type StreamBi interface {
	Recv() (GrpcDto, error)
	Send(msg []byte) error
}

type StreamClient interface {
	Recv() (GrpcDto, error)
}
