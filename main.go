package main

import (
	"dominus/app/interfaces/rest/inbound"
	"dominus/app/interfaces/rest/middlewares"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func restServer() {

	mid := middlewares.NewMiddleware()
	router := router.New()
	inbound.NewRestApi(router)


	
	log.Fatal(fasthttp.ListenAndServe(":8000", mid.APImiddleware(router.Handler)))
}

func main() {

}