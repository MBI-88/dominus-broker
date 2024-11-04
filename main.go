package main

import (
	"context"
	"dominus/app/domain/config"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"dominus/app/domain/topic"
	"dominus/app/interactors"
	"dominus/app/interfaces/database"
	f "dominus/app/interfaces/fasthttp/input"
	fm "dominus/app/interfaces/fasthttp/middlewares"
	"dominus/app/interfaces/fasthttp/output"
	g "dominus/app/interfaces/grpc/input"
	gf "dominus/app/interfaces/grpc/middlewares"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	mode       *bool
	system     chan os.Signal
	showBanner *bool
	banner     = `
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

// catches inital variables
func init() {
	mode = flag.Bool("prod", false, "set operation mode")
	showBanner = flag.Bool("banner", true, "show banner")

	flag.Usage = func() {
		info := fmt.Sprintf("[*] ***Dominus*** [*]\n")
		info += "mode: boolean\n"
		info += "banner: boolean\n"
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
	rls := rules.NewRule()
	restClient := output.NewRestClient()
	repo := database.NewRepository(env.Dsn, env.Database, rls, mongoClient)

	switch args {

	case "migrate":
		repo.Migrations(env.Collections)

	case "start":
		//Instances
		logs := event.NewLogs("./logs")
		topic := topic.NewTopic()
		events := event.NewEvent(
			repo,
			restClient,
			topic,
		)
		events.InitialLoad()
		inter := interactors.NewInteractor(repo, topic, events, logs)

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

		midF := fm.NewMiddleware()
		apiToken := fm.NewMiddlewareApiToken(env.ApiToken)
		allowedHost := fm.NewMiddlewareHot(env.Cidr)

		midF.AddMiddleware(apiToken, allowedHost)
		router := router.New()
		f.NewRestApi(router, inter)

		r := fasthttp.Server{
			Handler:                            midF.Middlewares(router.Handler),
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
		var opts []grpc.ServerOption

		midG := gf.NewMiddleware(env.ApiToken, logs)

		if errC == nil && errK == nil {
			creds, err := credentials.NewServerTLSFromFile(env.SslCert, env.KeyFile)
			if err != nil {
				cancel()
				panic(err)
			}
			opts = append(opts,
				grpc.Creds(creds),
				grpc.ChainUnaryInterceptor(
					auth.UnaryServerInterceptor(midG.ApiToken),
					logging.UnaryServerInterceptor(midG.LogErrors()),
				),
				grpc.ChainStreamInterceptor(
					auth.StreamServerInterceptor(midG.ApiToken),
					logging.StreamServerInterceptor(midG.LogErrors()),
				),
			)

		} else {
			opts = append(opts,
				grpc.ChainUnaryInterceptor(
					midG.UnaryLog,
					logging.UnaryServerInterceptor(midG.LogErrors()),
				),
				grpc.ChainStreamInterceptor(
					midG.StreamLog,
					logging.StreamServerInterceptor(midG.LogErrors()),
				),
			)
		}

		server := g.NewGrpcServe(opts, inter)
		listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", env.GrpcPort))

		go func(sr *grpc.Server, list net.Listener, cancel context.CancelFunc) {
			log.Fatal(sr.Serve(list))
			cancel()
		}(server, listener, cancel)

		//*********************************
		//**********Banner*****************
		//*********************************

		if errC == nil && errK == nil && *showBanner {
			fmt.Printf("%s 🚀 Rest: https://0.0.0.0:%d 🚀 Grpc: https://0.0.0.0:%d\n", banner, env.RestPort, env.GrpcPort)
		} else {
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

// Dominus entripoint
func main() {
	run()
}
