package interactors

import (
	"dominus/app/domain/entities"

	"github.com/valyala/fasthttp"
	jsoniter "github.com/json-iterator/go"
)


type interactor struct {
	m *entities.MessageRest // Rest struct
	// grpc struct
	parser jsoniter.API
}



func (i interactor) RestToGrp(ctx *fasthttp.RequestCtx) {

}


func (i interactor) RestToRest(ctx *fasthttp.RequestCtx) {

}


func (i interactor) GrpcToRest() {

}


func (i interactor) CreateTopic(ctx *fasthttp.RequestCtx) {

}


func (i interactor) UpdateTopic(ctx *fasthttp.RequestCtx) {

}

func (i interactor) GetTopic(ctx *fasthttp.RequestCtx) {


}




type InteractorInt interface {
	// RestToGrp sends information from  rest protocol receiver to  grpc protocol client
	RestToGrp(ctx *fasthttp.RequestCtx) 
	// RestToRest sends information from rest protocol receiver to rest protocol client
	RestToRest(ctx *fasthttp.RequestCtx)
	// GrpcToRest sends information from grpc protocol receiver to rest protocol client
	GrpcToRest()
	// UpdateNode updates topics and subcribers in the node
	UpdateTopic(ctx *fasthttp.RequestCtx)
	// Create a new topic in dominus
	CreateTopic(ctx *fasthttp.RequestCtx)
	// Return all topic in dominus
	GetTopic(ctx *fasthttp.RequestCtx)
}

// Create a new interactor instance
func NewInteractor() InteractorInt {
	return &interactor{m: new(entities.MessageRest), parser: jsoniter.ConfigCompatibleWithStandardLibrary}
}