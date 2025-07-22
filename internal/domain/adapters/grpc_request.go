package adapters

type IGrpcDto interface {
	GetPayload() []byte
	GetSubscribers() []string
}
