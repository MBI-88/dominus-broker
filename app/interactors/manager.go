package interactors

import (
	"dominus/app/domain/rules"
	"dominus/app/domain/topic"
	"dominus/app/interfaces/database"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type manager struct {
	p  jsoniter.API
	r  rules.RuleInt
	t  topic.TopicInt
	repo database.RepositoryInt
}


func (m manager) CreateTopic(ctx *fasthttp.RequestCtx) {
	temp := make(map[string][]string)
	message := make(map[string]string)
	body := ctx.Request.Body()

	if err := m.p.Unmarshal(body, temp); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	if err := m.r.ValidateTopic(temp); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	for k, v := range temp {
		m.t.CreateTopic(k, v)
	}

	message["message"] = "topic created successful!"
	b, _ := m.p.Marshal(message)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(b)
}

func (m manager) GetTopic(ctx *fasthttp.RequestCtx) {
	data := m.t.GetTopic() // cambiar por conexion con la bd
	body, err := m.p.Marshal(data)
	if err != nil {
		temp := make(map[string]string)
		temp["message"] = err.Error()
		body, _ := m.p.Marshal(temp)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Response.SetBody(body)
		return
	}
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(body)
}

func (m manager) DeleteTopic(ctx *fasthttp.RequestCtx) {
	temp := make(map[string]string)
	message := make(map[string]string)
	body := ctx.Request.Body()

	if err := m.p.Unmarshal(body, temp); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	if err := m.t.DeleteTopic(temp["topic"]); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	message["message"] = "topic deleted successful!"
	b, _ := m.p.Marshal(message)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(b)
}


type ManagerInt interface {
	//Create a new topic in dominus
	CreateTopic(ctx *fasthttp.RequestCtx)
	//Return all topic in dominus
	GetTopic(ctx *fasthttp.RequestCtx)
	//Delete a topic using a key selected
	DeleteTopic(ctx *fasthttp.RequestCtx)
}