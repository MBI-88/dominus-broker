package repos

import (
	"context"
	"dominus/app/domain/entities"
)

type RestContextInt interface {
	BodyParser(obj any) error
	Queries() map[string]string
	Params(key string) string
	QueryInt(key string) uint64
}

type RepositoryInt interface {
	Migrations()
	InsertObject(obj any, collection string) error
	FindObjects(collection string, objects any, filter any, page, size int) error
	FindObject(f any, collection string, object any) error
	DeleteObjects(f any, collection string) error
	Filter(filter any, object any, collection string) error
	CountPages(collection string) (int64, error)
	Stats() (any, error)
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
	ClientStream(urls []string, msg <-chan []byte, sig chan<- entities.Logs, tx chan<- struct{})
	ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, sig chan<- entities.Logs, closed <-chan struct{}, done chan<- struct{})
	BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- entities.Logs, tx chan<- struct{}, rx <-chan struct{}, done chan<- struct{})
}