package consumer

type ConsumerDto interface {
	GetMessageId() string
	GetWorkerId() string
	GetGroupId() string
}
