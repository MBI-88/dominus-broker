package dto

type BrokerRequestDto interface {
	GetPayload() []byte
	GetSubscribers() []string
}
