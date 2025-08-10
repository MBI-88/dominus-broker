package adapters

type GrpcResponse interface {
	GetMessage() string
	GetStatus() int64
}
