package input

import (
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	fastHttpSwagger "github.com/swaggo/fasthttp-swagger"
)

type swagger struct {
	r *router.Router
}


func (s *swagger) path() {
	s.r.GET("/swagger/*any", func(ctx *fasthttp.RequestCtx) {
		fastHttpSwagger.WrapHandler(fastHttpSwagger.InstanceName("swagger"))(ctx)
	})
}

func NewSwagger(r *router.Router ) {
	swg := &swagger{
		r: r,
	}
	swg.path()
}