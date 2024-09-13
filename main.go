package main

import (
	"dominus/app/domain/config"
	"dominus/app/interfaces/rest/inbound"
	"dominus/app/interfaces/rest/middlewares"
	"log"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)



func restDominus(mode bool) {
	// creation 
	conf := config.NewRestConfig().GetEnvVar(mode)
	mid := middlewares.NewMiddleware()
	router := router.New()
 	inbound.NewRestApi(router)

	// set options
	mid.SetApiToken(conf.ApiToken, conf.Cidr)

	dominus := fasthttp.Server{
		Handler:  mid.Middlewares(router.Handler),
		Name: "Dominus",
		ReadTimeout: time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
		IdleTimeout: time.Duration(conf.IdleTimeout),
		MaxConnsPerIP: conf.MaxConnsPerIp,
		MaxRequestsPerConn: conf.MaxRequestPerConn,
		MaxRequestBodySize: conf.MaxRequestBodySize,
		ReduceMemoryUsage: conf.ReduceMemoryUsage,
		DisablePreParseMultipartForm: conf.DisablePreparseMultipartForm,
		DisableHeaderNamesNormalizing: conf.DisableHeaderNamesNormalizing,
		SleepWhenConcurrencyLimitsExceeded: time.Duration(conf.SleepWhenConcurrencyLimitExcedeed),
		NoDefaultDate: conf.NoDefaultDate,
		KeepHijackedConns: conf.KeepHijackedConns,
		CloseOnShutdown: conf.CloseOnShutdown,
		StreamRequestBody: conf.StreamRequestBody,
	}


	// run server
	log.Fatal(dominus.ListenAndServe(":8000"))
}


func grpServer() {

}

func main() {

   go restDominus(false)
   go grpServer()

}