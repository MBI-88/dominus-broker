package ack

type AskDto interface {
	GetMessageId() string
	GetWorkerId() string
	GetGroupId() string
}
