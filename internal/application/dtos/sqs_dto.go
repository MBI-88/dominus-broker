package dtos

type ProducerDto interface {
	GetPayload() []byte
}

type ConsumerDto interface {
	GetId() string
	GetWorker() string
	GetGroupId() string
}
