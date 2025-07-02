package input

import (
	mg "dominus-project/internal/interactors/manager"
	"dominus-project/internal/interfaces/fasthttp/dto"
	"dominus-project/internal/interfaces/fasthttp/errors"

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
// @Accept json
// @Produce json
// @Param topic body dto.SwaggerTopic true "Topic Info"
// @Success 204 "No Content"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /topics [post]
func (m *manager) addTopic(c *fasthttp.RequestCtx) {
	ctx := dto.NewRestContext(c)
	if err := m.uc.AddTopic(ctx); err != nil {
		b := errors.SetErrorResponse(c, "application/json", fasthttp.StatusNotAcceptable, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

// @Tags Manager
// @Description <h3>update  subscribers</h3>
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param name path string true "topic name"
// @Param subscribers body []string true "Subscribers"
// @Success 204 "No Content"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /subscribers/{name} [patch]
func (m *manager) updateSubscribers(c *fasthttp.RequestCtx) {
	ctx := dto.NewRestContext(c)
	if err := m.uc.UpdateSubscribers(ctx); err != nil {
		b := errors.SetErrorResponse(c, "application/json", fasthttp.StatusNotAcceptable, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}

// @Tags Manager
// @Description <h3>delete a topic</h3>
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param name path string true "topic name"
// @Success 204 "No Content"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /topics/{name} [delete]
func (m *manager) deleteTopic(c *fasthttp.RequestCtx) {
	ctx := dto.NewRestContext(c)
	if err := m.uc.DeleteTopic(ctx); err != nil {
		b := errors.SetErrorResponse(c, "application/json", fasthttp.StatusNotAcceptable, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusNoContent)
}


// @Tags Manager
// @Description <h3>get topics info</h3>
// @Security ApiKeyAuth
// @Produce json
// @Success 200  {object} map[string]any "Topic response"
// @Failure 406 {object} map[string]string "Response body {message: error}"
// @Router /topics [get]
func (m *manager) getTopicsInfo(c *fasthttp.RequestCtx) {
	resp := m.uc.GetQueueInfo() 
	c.Response.Header.Set("Content-Type", "application/json")
	c.Response.Header.SetStatusCode(fasthttp.StatusOK)
	body, err := m.js.Marshal(resp)
	if err != nil {
		b := errors.SetErrorResponse(c, "application/json", fasthttp.StatusNotAcceptable, err.Error())
		c.Response.SetBody(b)
		return
	}
	c.Response.SetBody(body)

}

func (m *manager) path() {
	m.router.POST("/topics", m.addTopic)
	m.router.PATCH("/subscribers/{name}", m.updateSubscribers)
	m.router.DELETE("/topics/{name}", m.deleteTopic)
	m.router.GET("/topics", m.getTopicsInfo)
}

func NewManagerAPI(r *router.Router, uc mg.ManagerInt) {
	mg := &manager{
		router: r,
		uc:     uc,
		js: jsoniter.ConfigCompatibleWithStandardLibrary,
	}
	mg.path()
}
