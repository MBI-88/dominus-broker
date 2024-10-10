package main

import (
	"context"
	"dominus/app/domain/config"
	"dominus/app/domain/event"
	"dominus/app/domain/topic"
	"dominus/app/interactors"
	"dominus/app/interfaces/clients"
	"dominus/app/interfaces/database"
	"dominus/app/interfaces/rest/input"
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
	banner = `
==========================================================	
██████    ██████  ███    ██  ██ ███    ██ ██    ██ ███████
██   ██  ██    ██ ████  ████ ██ ████   ██ ██    ██ ██
██   ██  ██    ██ ██ ████ ██ ██ ██ ██  ██ ██    ██ ███████
██   ██  ██    ██ ██  ██  ██ ██ ██  ██ ██ ██    ██      ██
██████    ██████  ██      ██ ██ ██   ████  ██████  ███████
==========================================================    
   
👉 Github: https://github.com/MBI-88
🔧 Press CTRL+C to terminate the server

Dominus server is running on`
)

//catches inital variables
func init() {
	mode = flag.Bool("prod", false, "set operation mode")

	flag.Usage = func() {
		info := fmt.Sprintf("[*] ***Dominus*** [*]\n")
		info += "mode: boolean\n"
		info += "args: migrate|start"
		fmt.Fprintf(os.Stderr, "%s\n", info)
		flag.PrintDefaults()
	}
}

func run() {
	//Receives commands from cli
	flag.Parse()
	args := flag.Arg(0)

	//Instances
	settings := config.NewConfig()
	env := settings.GetEnvVar(*mode)
	mongoConfig := database.NewMongoConfig()
	mongoClient := mongoConfig.CreateClient(env.Dsn)

	client := clients.NewClient(env.Dsn, env.Database, mongoClient)

	switch args {

	case "migrate":
		repo := client.MongoClient()
		repo.Migrations(env.Collections)

	case "start":
		//Instances
		logs := event.NewLogs("./logs")
		topic := topic.NewTopic()
		events := event.NewEvent(
			client.MongoClient(),
			client.RestClient(),
			topic,
		)
		events.InitialLoad()
		inter := interactors.NewInteractor(client, topic, events, logs)

		// Signals
		system = make(chan os.Signal, 1)
		signal.Notify(system, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(system)

		//context
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		//**************************************
		//***********Rest Server****************
		//**************************************

		mid := middlewares.NewMiddleware()
		apiToken := middlewares.NewMiddlewareApiToken(env.ApiToken)
		allowedHost := middlewares.NewMiddlewareHot(env.Cidr)


		mid.AddMiddleware(apiToken,allowedHost)
		router := router.New()
		input.NewRestApi(router, inter)

		r := fasthttp.Server{
			Handler:                            mid.Middlewares(router.Handler),
			Name:                               "Dominus",
			ReadTimeout:                        time.Duration(env.ReadTimeout) * time.Second,
			WriteTimeout:                       time.Duration(env.WriteTimeout) * time.Second,
			IdleTimeout:                        time.Duration(env.IdleTimeout),
			MaxConnsPerIP:                      env.MaxConnsPerIp,
			MaxRequestsPerConn:                 env.MaxRequestPerConn,
			MaxRequestBodySize:                 env.MaxRequestBodySize,
			ReduceMemoryUsage:                  env.ReduceMemoryUsage,
			DisablePreParseMultipartForm:       env.DisablePreparseMultipartForm,
			DisableHeaderNamesNormalizing:      env.DisableHeaderNamesNormalizing,
			SleepWhenConcurrencyLimitsExceeded: time.Duration(env.SleepWhenConcurrencyLimitExcedeed),
			NoDefaultDate:                      env.NoDefaultDate,
			KeepHijackedConns:                  env.KeepHijackedConns,
			CloseOnShutdown:                    env.CloseOnShutdown,
			StreamRequestBody:                  env.StreamRequestBody,
			Logger:                             logs,
		}	

		_, errC := os.Stat(env.SslCert)
		_, errK := os.Stat(env.KeyFile)
		

		if errC == nil && errK == nil {
			go func(port uint16, cert, key string, cancel context.CancelFunc) {
				log.Fatal(r.ListenAndServeTLS(fmt.Sprintf(":%d", port), cert, key))
				cancel()
			}(env.RestPort, env.SslCert, env.KeyFile, cancel)
		} else {
			go func(port uint16, cancel context.CancelFunc) {
				log.Fatal(r.ListenAndServe(fmt.Sprintf(":%d", port)))
				cancel()
			}(env.RestPort, cancel)
		}

		//********************************
		//*********Grpc server************
		//********************************

		
		//*********************************
		//**********Banner*****************
		//*********************************

		if errC == nil && errK == nil {
			fmt.Printf("%s 🚀 Rest: https://0.0.0.0:%d 🚀 Grpc: https://0.0.0.0:%d\n", banner, env.RestPort, env.GrpcPort)
		}else {
			fmt.Printf("%s 🚀 Rest: http://0.0.0.0:%d 🚀 Grpc: http://0.0.0.0:%d\n", banner, env.RestPort, env.GrpcPort)
		}

		//*********************************
		//*********Shutdown servers********
		//*********************************

		// Wait for a signal
		select {
		case <-system:
			break
		case <-ctx.Done():
			break
		}

		if err := r.Shutdown(); err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println("No option selected")
	}
}



//Dominus entripoint
func main() {
	run()
}
