package input

import (
	"github.com/fasthttp/router"
	fastHttpSwagger "github.com/swaggo/fasthttp-swagger"
	"github.com/valyala/fasthttp"
)

type swagger struct {
	router *router.Router
}

func NewSwaggerAPI(r *router.Router) {
	swg := &swagger{router: r}
	swg.path()
}

func (s *swagger) path() {
	s.router.GET("/swagger/{*}", func(ctx *fasthttp.RequestCtx) {
		fastHttpSwagger.WrapHandler(fastHttpSwagger.InstanceName("swagger"))(ctx)
	})
}
