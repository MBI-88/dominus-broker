package adapters

type GrpcDto interface {
	GetPayload() []byte
	GetSubscribers() []string
}
