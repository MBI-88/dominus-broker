package dtos

type ProducerDto interface {
	GetPayload() []byte
}

type ConsumerDto interface {
	GetMessageId() string
	GetWorkerId() string
	GetGroupId() string
}
