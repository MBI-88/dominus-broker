package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/rules"
	"dominus/app/domain/topic"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type interactor struct {
	m *entities.Message // Rest struct
	// grpc struct
	p jsoniter.API
	r rules.RuleInt
	t topic.TopicInt
}

func (i interactor) RestToGrp(ctx *fasthttp.RequestCtx) {

}

func (i interactor) RestToRest(ctx *fasthttp.RequestCtx) {

}

func (i interactor) GrpcToRest() {

}

func (i interactor) CreateTopic(ctx *fasthttp.RequestCtx) {
	temp := make(map[string][]string)
	message := make(map[string]string)
	body := ctx.Request.Body()

	if err := i.p.Unmarshal(body, temp); err != nil {
		message["message"] = err.Error()
		b, _ := i.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	if err := i.r.ValidateTopic(temp); err != nil {
		message["message"] = err.Error()
		b, _ := i.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	for k, v := range temp {
		i.t.CreateTopic(k, v)
	}

	message["message"] = "topic created successful!"
	b, _ := i.p.Marshal(message)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(b)
}

func (i interactor) GetTopic(ctx *fasthttp.RequestCtx) {
	data := i.t.GetTopic() // cambiar por conexion con la bd
	body, err := i.p.Marshal(data)
	if err != nil {
		temp := make(map[string]string)
		temp["message"] = err.Error()
		body, _ := i.p.Marshal(temp)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Response.SetBody(body)
		return
	}
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(body)
}

func (i interactor) DeleteTopic(ctx *fasthttp.RequestCtx) {
	temp := make(map[string]string)
	message := make(map[string]string)
	body := ctx.Request.Body()

	if err := i.p.Unmarshal(body, temp); err != nil {
		message["message"] = err.Error()
		b, _ := i.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	if err := i.t.DeleteTopic(temp["topic"]); err != nil {
		message["message"] = err.Error()
		b, _ := i.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	message["message"] = "topic deleted successful!"
	b, _ := i.p.Marshal(message)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(b)
}

type InteractorInt interface {
	// RestToGrp sends information from  rest protocol receiver to  grpc protocol client
	RestToGrp(ctx *fasthttp.RequestCtx)
	// RestToRest sends information from rest protocol receiver to rest protocol client
	RestToRest(ctx *fasthttp.RequestCtx)
	// GrpcToRest sends information from grpc protocol receiver to rest protocol client
	GrpcToRest()
	// Create a new topic in dominus
	CreateTopic(ctx *fasthttp.RequestCtx)
	// Return all topic in dominus
	GetTopic(ctx *fasthttp.RequestCtx)
	// Delete a topic using a key selected
	DeleteTopic(ctx *fasthttp.RequestCtx)
}

// Create a new interactor instance
func NewInteractor() InteractorInt {
	return &interactor{
		m: new(entities.Message),
		p: jsoniter.ConfigCompatibleWithStandardLibrary,
		r: rules.NewRule(),
		t: topic.NewTopic()}
}
