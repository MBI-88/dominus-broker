package repos


type StreamBiInt interface {
	Recv() (GrpRequestMessageInt, error)
	Send(msg []byte) error
}

type StreamClientInt interface {
	Recv() (GrpRequestMessageInt, error)
}