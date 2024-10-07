package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/event"
	"dominus/app/domain/topic"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type connection struct {
	m *entities.Message // Rest struct
	// grpc struct
	p  jsoniter.API
	t  topic.TopicInt
	ev event.EventsInt
	// grpc client
	// resp client
}

func (c connection) RestToGrpc(ctx *fasthttp.RequestCtx) {

}

func (c connection) RestToRest(ctx *fasthttp.RequestCtx) {
	
}

func (i connection) GrpcToRest() {

}


type ConnectionInt interface {
	//RestToGrp sends information from  rest protocol receiver to  grpc protocol client
	//
	//Parameters
	//
	//-> ctx: fasthttp context
	RestToGrpc(ctx *fasthttp.RequestCtx)
	//RestToRest sends information from rest protocol receiver to rest protocol client
	//
	//Parameters
	//
	//-> ctx: fasthttp context
	RestToRest(ctx *fasthttp.RequestCtx)
	//GrpcToRest sends information from grpc protocol receiver to rest protocol client
	GrpcToRest()
}

