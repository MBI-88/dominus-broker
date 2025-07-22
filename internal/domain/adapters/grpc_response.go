package adapters

type IGrpcResponse interface {
	GetMessage() string
	GetStatus() uint32

}
