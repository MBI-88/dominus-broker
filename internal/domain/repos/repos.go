package repos

import (
	"context"
)

type RestContextInt interface {
	BodyParser(obj any) error
	Param(key string) string
}

type GrpRequestMessageInt interface {
	Descriptor() ([]byte, []int)
	GetPayload() []byte
	GetSubscribers() []string
	Reset()
	String() string
	Validate() error
	ValidateAll() error
}

type StreamClientInt interface {
	Recv() (GrpRequestMessageInt, error)
}

type StreamServerInt interface {
	Send(payload []byte) error
	Context() context.Context
}

type StreamBiInt interface {
	Recv() (GrpRequestMessageInt, error)
	Send(msg []byte) error
}

type GrpResponseInt interface {
	Descriptor() ([]byte, []int)
	GetMessage() string
	GetStatus() uint32
	ProtoMessage()
	Reset()
	String() string
	Validate() error
	ValidateAll() error
}

type GrpClientInt interface {
	Simple(url string, msg []byte) (GrpResponseInt, error)
	ClientStream(urls []string, msg <-chan []byte, ctx context.Context)
	ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, ctx context.Context, done chan<- struct{})
	BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- error, tx chan<- struct{}, ctx context.Context, done chan<- struct{})
}

type LogsInt interface {
	WriteLog(op, dsc string)
	Printf(format string, args ...any)
	GetLogs() ([]string, error)
}


type TopicInt interface {
	SetMessage(data []byte)
	GetMessage() []byte
	SetSubscribers(sb []string)
	GetSubscribers() []string
	GetName() string
	SetName(name string)
}