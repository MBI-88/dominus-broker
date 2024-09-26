package interactors

import (
	"dominus/app/domain/entities"
	"dominus/app/domain/rules"
	"dominus/app/domain/topic"
	"dominus/app/interfaces/database"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type manager struct {
	p          jsoniter.API
	r          rules.RuleInt
	t          topic.TopicInt
	repo       database.RepositoryInt
	tdb        *entities.TopicDB
	collection string
}

func (m manager) CreateTopic(ctx *fasthttp.RequestCtx) {
	var (
		message = make(map[string]any)
		body    = ctx.Request.Body()
	)

	if err := m.p.Unmarshal(body, m.tdb); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	if err := m.r.ValidateStruct(m.tdb); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	m.t.CreateTopic(m.tdb.Topic, m.tdb.Subscribers)

	if _, err := m.repo.InsertObject(m.tdb, m.collection); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Response.SetBody(b)
		return
	}

	message["message"] = "topic created successful!"
	b, _ := m.p.Marshal(message)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(b)
}

func (m manager) GetTopic(ctx *fasthttp.RequestCtx) {
	var (
		topics  []entities.TopicDB
		message = make(map[string]any)
	)

	if err := m.repo.FindObjects(m.collection, &topics, bson.D{}); err != nil {
		message["message"] = err.Error()
		body, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Response.SetBody(body)
		return
	}

	message["topics"] = topics

	body, err := m.p.Marshal(message)
	if err != nil {
		delete(message, "topics")
		message["message"] = err.Error()
		body, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.Response.SetBody(body)
		return
	}
	ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.SetBody(body)
}

func (m manager) UpdateTopic(ctx *fasthttp.RequestCtx) {
	message := make(map[string]string)
	body := ctx.Request.Body()

	if err := m.p.Unmarshal(body, m.tdb); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	if err := m.r.ValidateStruct(m.tdb); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	m.t.CreateTopic(m.tdb.Topic, m.tdb.Subscribers)
	filter := primitive.D{primitive.E{Key: "topic", Value: m.tdb.Topic}}
	update := primitive.D{primitive.E{Key: "$set", Value: m.tdb}}

	if _, err := m.repo.UpdateObject(filter, update, m.collection); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	message["message"] = "topic updated successful!"
	b, _ := m.p.Marshal(message)
	ctx.Response.Header.SetStatusCode(fasthttp.StatusAccepted)
	ctx.Response.SetBody(b)
}

func (m manager) DeleteTopic(ctx *fasthttp.RequestCtx) {
	message := make(map[string]string)
	body := ctx.Request.Body()

	if err := m.p.Unmarshal(body, m.tdb); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	if err := m.r.ValidateStruct(m.tdb); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	if err := m.t.DeleteTopic(m.tdb.Topic); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.Response.SetBody(b)
		return
	}

	filter := primitive.D{primitive.E{Key: "topic", Value: m.tdb.Topic}}

	if _, err := m.repo.DeleteObject(filter, m.collection); err != nil {
		message["message"] = err.Error()
		b, _ := m.p.Marshal(message)
		ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
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
	//Update a topic in dominus
	UpdateTopic(ctx *fasthttp.RequestCtx)
}
