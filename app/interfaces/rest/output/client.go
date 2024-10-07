package output

import (
	"time"

	"github.com/valyala/fasthttp"
)




type restclient struct {
	c *fasthttp.Client
}

func (r restclient) DoJsonRequest(sub string, payload []byte) error {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(sub)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.SetBody(payload)
	
	if err := r.c.Do(req,resp); err != nil {
		return err
	}

	return nil
}



type RestClientInt interface {
	//Rest client
	//
	//Parameters
	//
	//-> ctx: context
	// 
	//-> sub: subcriber
	//
	//-> payload: message
	DoJsonRequest(sub string, payload []byte) error
}

func NewRestClient() RestClientInt {
	return &restclient{
		c: &fasthttp.Client{
			Name: "Dominus",
			MaxConnsPerHost: 1,
			MaxIdleConnDuration: 2 * time.Second,
			MaxConnDuration: 2 * time.Second,
			MaxIdemponentCallAttempts: 1,
		},
	}
}