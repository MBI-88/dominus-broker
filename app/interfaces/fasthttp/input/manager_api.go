package input

import (
	"dominus-project/app/interactors"

	"github.com/fasthttp/router"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type manager struct {
	router  *router.Router
	js      jsoniter.API
	service interactors.ManagerInt
}

func (m *manager) addTopic(c *fasthttp.RequestCtx) {
	ctx := NewRestContext(c)
	if err := m.service.AddTopic(ctx); err != nil {
		b := setErrorResponse(c, "application/json", fasthttp.StatusInternalServerError, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

func (m *manager) updatePartition(c *fasthttp.RequestCtx) {
	ctx := NewRestContext(c)
	if err := m.service.UpdatePartition(ctx); err != nil {
		b := setErrorResponse(c, "application/json", fasthttp.StatusInternalServerError, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

func (m *manager) deleteTopic(c *fasthttp.RequestCtx) {
	ctx := NewRestContext(c)
	if err := m.service.DeleteTopic(ctx); err != nil {
		b := setErrorResponse(c, "application/json", fasthttp.StatusInternalServerError, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}



func (m *manager) path() {
	m.router.POST("/add-topic", m.addTopic)
	m.router.PUT("/update-partition/{name}", m.updatePartition)
	m.router.DELETE("/delete-topic/{name}", m.deleteTopic)
}

func NewManagerAPI(r *router.Router, srv interactors.ManagerInt) {
	mg := &manager{
		router: r,
		service: srv,
	}
	mg.path()
}
