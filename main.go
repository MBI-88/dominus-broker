package main

import (
	"context"
	"dominus/app/domain/config"
	"dominus/app/interfaces/rest/inbound"
	"dominus/app/interfaces/rest/middlewares"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	mode   *bool
	system chan os.Signal
)

//Starts rest service
func runRestServer(mode bool, cancel context.CancelFunc) *fasthttp.Server {
	// creation
	conf := config.NewConfig().GetEnvVar(mode)
	mid := middlewares.NewMiddleware()
	router := router.New()
	inbound.NewRestApi(router)

	// set options
	mid.SetApiToken(conf.ApiToken, conf.Cidr)

	s := fasthttp.Server{
		Handler:                            mid.Middlewares(router.Handler),
		Name:                               "Dominus",
		ReadTimeout:                        time.Duration(conf.ReadTimeout) * time.Second,
		WriteTimeout:                       time.Duration(conf.WriteTimeout) * time.Second,
		IdleTimeout:                        time.Duration(conf.IdleTimeout),
		MaxConnsPerIP:                      conf.MaxConnsPerIp,
		MaxRequestsPerConn:                 conf.MaxRequestPerConn,
		MaxRequestBodySize:                 conf.MaxRequestBodySize,
		ReduceMemoryUsage:                  conf.ReduceMemoryUsage,
		DisablePreParseMultipartForm:       conf.DisablePreparseMultipartForm,
		DisableHeaderNamesNormalizing:      conf.DisableHeaderNamesNormalizing,
		SleepWhenConcurrencyLimitsExceeded: time.Duration(conf.SleepWhenConcurrencyLimitExcedeed),
		NoDefaultDate:                      conf.NoDefaultDate,
		KeepHijackedConns:                  conf.KeepHijackedConns,
		CloseOnShutdown:                    conf.CloseOnShutdown,
		StreamRequestBody:                  conf.StreamRequestBody,
	}

	fmt.Printf("[*] Rest service running on 0.0.0.0:%d\n", conf.Port)

	// run server
	if conf.SslCert != "" && conf.KeyFile != "" {
		go func(port uint16, cert, key string, cancel context.CancelFunc) {
			log.Fatal(s.ListenAndServeTLS(fmt.Sprintf(":%d", port), cert, key))
			cancel()
		}(conf.Port, conf.SslCert, conf.KeyFile, cancel)
	} else {
		go func(port uint16, cancel context.CancelFunc) {
			log.Fatal(s.ListenAndServe(fmt.Sprintf(":%d", port)))
			cancel()
		}(conf.Port, cancel)
	}

	return &s
}

//Starts gRPC service
func runGrpServer(mode bool, cancel context.CancelFunc) {
}

//catches inital variables
func init() {
	mode = flag.Bool("mode", false, "set operation mode")

	flag.Usage = func() {
		info := fmt.Sprintf("[*] ***Dominus*** [*]")
		info = "\nmode: boolean\n"

		fmt.Fprintf(os.Stderr, "%s\n", info)
		flag.PrintDefaults()
	}
}

//Dominus entripoint
func main() {
	// Receives commands from cli
	flag.Parse()

	//signal
	system = make(chan os.Signal, 1)
	signal.Notify(system, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(system)

	//context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	//servers
	r := runRestServer(*mode, cancel)
	runGrpServer(*mode, cancel)

	// Wait for a signal
	select {
	case <-system:
		break
	case <-ctx.Done():
		break

	}

	// Shutdown servers
	if err := r.Shutdown(); err != nil {
		log.Fatal(err)
	}


}
