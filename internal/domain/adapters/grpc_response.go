package adapters

type GrpcResponse interface {
	GetMessage() string
	GetStatus() uint32
}
