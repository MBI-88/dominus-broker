package main

import (
	"context"
	"dominus/app/domain/config"
	"dominus/app/domain/event"
	"dominus/app/domain/rules"
	"dominus/app/interactors"
	"dominus/app/interfaces/database"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	fi "dominus/app/interfaces/fasthttp/input"
	fm "dominus/app/interfaces/fasthttp/middlewares"
	gi "dominus/app/interfaces/grpc/input"
	gm "dominus/app/interfaces/grpc/middlewares"
	gt "dominus/app/interfaces/grpc/output"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"time"

	grpcmetrics "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

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
	repo := database.NewRepository(env.Dsn, env.Database, env.Collections, rls, mongoClient)

	switch args {
	case "migrate":
		repo.Migrations()
	case "start":
		//Instances
		logs := event.NewLogs("./logs")

		inter := interactors.NewInteractor(repo, logs)

		// Signals
		system = make(chan os.Signal, 1)
		signal.Notify(system, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(system)

		//context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		//**************************************
		//***********Monitoring*****************
		//**************************************

		grpcmetrics := grpcmetrics.NewServerMetrics(
			grpcmetrics.WithServerHandlingTimeHistogram(
				grpcmetrics.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
			),
		)

		reg := prometheus.NewRegistry()
		reg.MustRegister(grpcmetrics)

		//**************************************
		//***********Rest Server****************
		//**************************************

		midF := fm.NewMiddleware()
		apiToken := fm.NewMiddlewareApiToken(env.ApiToken)
		allowedHost := fm.NewMiddlewareHost(env.Cidr)

		midF.AddMiddleware(apiToken, allowedHost)
		router := router.New()
		fi.NewRestController(router, inter)
		fi.NewMonitor(router, reg)

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
		_, errCa := os.Stat(env.SslCaCert)

		if errC == nil && errK == nil {
			go func(port uint64, cert, key string, cancel context.CancelFunc) {
				log.Fatal(r.ListenAndServeTLS(fmt.Sprintf("0.0.0.0:%d", port), cert, key))
				cancel()
			}(env.RestPort, env.SslCert, env.KeyFile, cancel)
		} else {
			go func(port uint64, cancel context.CancelFunc) {
				log.Fatal(r.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port)))
				cancel()
			}(env.RestPort, cancel)
		}

		//********************************
		//*********Grpc server************
		//********************************
		var (
			optsS []grpc.ServerOption
			optsD []grpc.DialOption
		)

		midGs := gm.NewMiddleware(env.ConnectionKey, logs)
		midGc := gm.NewInterceptor(env.ConnectionKey)

		if errC == nil && errK == nil && errCa == nil {
			credsS, err := credentials.NewServerTLSFromFile(env.SslCert, env.KeyFile)
			if err != nil {
				cancel()
				panic(err)
			}
			optsS = append(optsS, grpc.Creds(credsS))

			credsD, err := credentials.NewClientTLSFromFile(env.SslCaCert, "dominus.com")
			if err != nil {
				cancel()
				panic(err)
			}
			optsD = append(optsD, grpc.WithTransportCredentials(credsD))
		} else {
			optsD = append(optsD, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}

		optsS = append(optsS,
			grpc.ChainUnaryInterceptor(
				otelgrpc.UnaryServerInterceptor(),
				grpcmetrics.UnaryServerInterceptor(),
				auth.UnaryServerInterceptor(midGs.ApiToken),
				logging.UnaryServerInterceptor(midGs.LogErrors()),
			),
			grpc.ChainStreamInterceptor(
				otelgrpc.StreamServerInterceptor(),
				grpcmetrics.StreamServerInterceptor(),
				auth.StreamServerInterceptor(midGs.ApiToken),
				logging.StreamServerInterceptor(midGs.LogErrors()),
			),
		)

		optsD = append(optsD,
			grpc.WithUnaryInterceptor(midGc.UnaryAuthInterceptor),
			grpc.WithStreamInterceptor(midGc.StreamAuthInterceptor),
			grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
			grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`),
		)

		gclient := gt.NewGrpClient(optsD)
		inter = inter.Set(gclient)
		srG := gi.NewGrpcController(optsS, inter)
		listener, _ := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", env.GrpcPort))

		go func(sr *grpc.Server, list net.Listener, cancel context.CancelFunc) {
			log.Fatal(sr.Serve(list))
			cancel()
		}(srG, listener, cancel)

		//*********************************
		//**********Banner*****************
		//*********************************

		if errC == nil && errK == nil && *showBanner {
			fmt.Printf("%s Rest: https://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", banner, env.RestPort, env.GrpcPort)
		} else {
			fmt.Printf("%s Rest: http://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", banner, env.RestPort, env.GrpcPort)
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
		srG.GracefulStop()

	default:
		fmt.Println("No option selected")
	}
}

// Dominus entripoint
func main() {
	run()
}
