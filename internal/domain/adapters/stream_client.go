package adapters

type IStreamBi interface {
	Recv() (IGrpcDto, error)
	Send(msg []byte) error
}

type IStreamClient interface {
	Recv() (IGrpcDto, error)
}
