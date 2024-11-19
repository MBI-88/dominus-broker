package interactors

import (
	"context"
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"mime/multipart"
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
	}
}

/*
func (i *interactor) NewManager() ManagerInt {
	return &manager{
		r:          rules.NewRule(),
		repo:       i.repo,
		collection: "topic",
		lg:         i.log,
	}
}
**/

func (i *interactor) Set(c GrpClientInt) InteractorInt {
	i.gclient = c
	return i
}

type InteractorInt interface {
	NewConnection() ConnectionInt
	//NewManager() ManagerInt
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
	//BodyParser parses context data to struct
	//
	//Parameters
	//
	//-> obj: struct to fill
	BodyParser(obj any) error
	//FormFile gives a multipart body
	//
	//Parameters
	//
	//-> key: the field to find
	//
	//Returns
	//
	//-> multipart.FileHeader
	//
	//-> error
	FormFile(key string) (*multipart.FileHeader, error)
	//FormValue gives a array byte of the key selected
	//
	//Parameters
	//
	//-> key: field to find
	//
	//Returns
	//
	//-> data: data array byte
	FormValue(key string) []byte
}

type RepositoryInt interface {
	//Create an object in the database
	//
	//Parameters
	//
	//-> obj: the object to be updated
	//
	//-> collection: collection name to use
	InsertObject(obj any, collection string) error
	//Find objects in the database
	//
	//Parameters
	//
	//-> collection: collection name to use
	//
	//-> objects: array object to fill
	//
	//-> filter: the filter to find objects
	//
	//-> page: page seletected
	//
	//-> size: total elements in the page
	FindObjects(collection string, objects any, filter any, page, size int) error
	//Find and object in the database
	//
	//Parameters
	//
	//-> f: filter to match with objects
	//
	//-> collection: collection name to use
	//
	//-> object: the object to fill
	FindObject(f any, collection string, object any) error
	//Delete objects in the database
	//
	//Parameters
	//
	//-> f: filter to find objects
	//
	//-> collection: collection name to use
	DeleteObjects(f any, collection string) error
}

type RestClientInt interface {
	//Rest client
	//
	//Parameters
	//
	//-> ctx: context
	//
	//-> sub: subcriber
	//
	//-> payload: message
	DoJsonRequest(sub string, payload []byte) error
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
	ClientStream(urls []string, msg <-chan []byte, sig chan<- *entities.Logs)
	ServerStream(urls []string, initalMsg []byte, msg chan<- []byte, sig chan<- *entities.Logs)
	BidirectionalStream(urls []string, provMsg <-chan []byte, subMsg chan<- []byte, errMsg chan<- *entities.Logs)
}
