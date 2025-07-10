package repos


type IStreamBi interface {
	Recv() (IGrpRequestMessage, error)
	Send(msg []byte) error
}

type IStreamClient interface {
	Recv() (IGrpRequestMessage, error)
}