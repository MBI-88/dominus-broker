package interactors

import (
	"context"
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
)

type interactor struct {
	repo    RepositoryInt
	log     event.LogsInt
	gclient GrpClientInt
}

func (i *interactor) NewConnection() ConnectionInt {
	return &connection{
		r:  rules.NewRule(),
		lg: i.log,
		collection: "logs",
		client: i.gclient,
		repo: i.repo,
	}
}


func (i *interactor) NewManager() ManagerInt {
	return &manager{
		rls:          rules.NewRule(),
		repo:       i.repo,
		collection: "logs",
		lg:         i.log,
	}
}


func (i *interactor) Set(c GrpClientInt) InteractorInt {
	i.gclient = c
	return i
}

type InteractorInt interface {
	NewConnection() ConnectionInt
	NewManager() ManagerInt
	Set(cl GrpClientInt) InteractorInt
}

// Create a new interactor instance
func NewInteractor(rp RepositoryInt, lg event.LogsInt) InteractorInt {
	return &interactor{
		repo: rp,
		log:  lg,
	}
}

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
	CountPages(collection string) (uint64, error)
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
	ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, sig chan<- entities.Logs, closed <-chan struct{}, done chan <-struct{})
	BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- entities.Logs, ctx context.Context)
}
