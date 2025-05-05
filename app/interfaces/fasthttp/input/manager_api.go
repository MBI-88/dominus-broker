package input

import (
	mg "dominus-project/app/interactors/manager"

	"github.com/fasthttp/router"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type manager struct {
	router *router.Router
	js     jsoniter.API
	uc     mg.ManagerInt
}

// @Tags Manager
// @Description <h3>add a new topic</h3>
// @Security ApiKeyAuth
// @Accept application/json
// @Param name body string true "topic name"
// @Param partitions body []string true "partitions"
// @Success 204 {object} map[string][]string "Success response"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /topic [post]
func (m *manager) addTopic(c *fasthttp.RequestCtx) {
	ctx := NewRestContext(c)
	if err := m.uc.AddTopic(ctx); err != nil {
		b := setErrorResponse(c, "application/json", fasthttp.StatusNotAcceptable, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

// @Tags Manager
// @Description <h3>update  a partition</h3>
// @Security ApiKeyAuth
// @Accept application/json
// @Param name path string true "topic name"
// @Param partitions body []string true "partitions"
// @Success 204 {object} map[string][]string "Success response"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /partition/{name} [put]
func (m *manager) updatePartition(c *fasthttp.RequestCtx) {
	ctx := NewRestContext(c)
	if err := m.uc.UpdatePartition(ctx); err != nil {
		b := setErrorResponse(c, "application/json", fasthttp.StatusNotAcceptable, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

// @Tags Manager
// @Description <h3>delete a topic</h3>
// @Security ApiKeyAuth
// @Accept application/json
// @Param name path string true "topic name"
// @Success 204 {object} map[string][]string "Success response"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /topic/{name} [delete]
func (m *manager) deleteTopic(c *fasthttp.RequestCtx) {
	ctx := NewRestContext(c)
	if err := m.uc.DeleteTopic(ctx); err != nil {
		b := setErrorResponse(c, "application/json", fasthttp.StatusNotAcceptable, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

func (m *manager) path() {
	m.router.POST("/topic", m.addTopic)
	m.router.PUT("/partition/{name}", m.updatePartition)
	m.router.DELETE("/topic/{name}", m.deleteTopic)
}

func NewManagerAPI(r *router.Router, uc mg.ManagerInt) {
	mg := &manager{
		router:  r,
		uc: uc,
	}
	mg.path()
}
