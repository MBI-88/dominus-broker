package main

import (
	"context"
	"dominus-project/config"
	"dominus-project/docs"
	"dominus-project/internal/domain/entities"
	grpcconn "dominus-project/internal/interactors/grpc_conn"
	"dominus-project/internal/interactors/manager"
	"dominus-project/internal/interactors/system"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	"dominus-project/internal/interfaces/events"
	fi "dominus-project/internal/interfaces/fasthttp/input"
	fm "dominus-project/internal/interfaces/fasthttp/middlewares"
	gi "dominus-project/internal/interfaces/grpc/input"
	gm "dominus-project/internal/interfaces/grpc/middlewares"
	gt "dominus-project/internal/interfaces/grpc/output"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/fasthttp/router"
	grpcmetrics "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	mode       *bool
	st         chan os.Signal
	showBanner *bool
	banner     = `
==========================================================	
██████    ██████  ███     ██ ██ ███    ██ ██    ██ ███████
██   ██  ██    ██ ████  ████ ██ ████   ██ ██    ██ ██
██   ██  ██    ██ ██ ████ ██ ██ ██ ██  ██ ██    ██ ███████
██   ██  ██    ██ ██  ██  ██ ██ ██  ██ ██ ██    ██      ██
██████    ██████  ██      ██ ██ ██   ████  ██████  ███████
==========================================================    
   
👉 Github: https://github.com/MBI-88
🔧 Press CTRL+C to terminate the server

dominus server is running on`
)

// catches inital variables
func init() {
	mode = flag.Bool("prod", false, "set operation mode")
	showBanner = flag.Bool("banner", true, "show banner")

	flag.Usage = func() {
		info := "[*] ***dominus-project*** [*]\n"
		info += "mode: boolean\n"
		info += "banner: boolean\n"
		fmt.Fprintf(os.Stderr, "%s\n", info)
		flag.PrintDefaults()
	}
}

// @title dominus-project server
// @description <h3>This server is a bidirectional queue using gRCP. Manager section</h3>
// @contact.name MBI
// @contact.email ingmbi8807@gmail.com
// @contact.url https://www.pr0c0d3.com/
// @version 1.1.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name x-api-key
// @host localhost:8000
// @BasePath /
func main() {
	//Receives commands from cli
	flag.Parse()

	//Instances
	settings := config.NewConfig()
	env := settings.GetEnvVar(*mode)

	logs := events.NewLogs(env.Logs)
	topics := entities.NewTopics(env.TopicLimit)
	docs.SwaggerInfo.Host = env.Host

	// Signals
	st = make(chan os.Signal, 1)
	signal.Notify(st, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(st)

	//context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//**************************************
	//***********Monitoring*****************
	//**************************************

	metricserver := grpcmetrics.NewServerMetrics(
		grpcmetrics.WithServerCounterOptions(grpcmetrics.WithConstLabels(prometheus.Labels{})),
		grpcmetrics.WithServerHandlingTimeHistogram(grpcmetrics.WithHistogramOpts(&prometheus.HistogramOpts{})),
	)
	metricclient := grpcmetrics.NewClientMetrics(
		grpcmetrics.WithClientCounterOptions(grpcmetrics.WithConstLabels(prometheus.Labels{})),
		grpcmetrics.WithClientHandlingTimeHistogram(grpcmetrics.WithHistogramOpts(&prometheus.HistogramOpts{})),
		grpcmetrics.WithClientStreamRecvHistogram(grpcmetrics.WithHistogramOpts(&prometheus.HistogramOpts{})),
	)

	reg := prometheus.NewRegistry()
	reg.MustRegister(metricserver, metricclient)

	//**************************************
	//***********Rest Server****************
	//**************************************

	midF := fm.NewMiddleware()
	apiToken := fm.NewMiddlewareApiToken(env.ApiToken)
	allowedOrings := fm.NewMiddlewareHost(env.AllowOrigins)

	// Interactors
	system := system.NewSystemService(logs)
	manager := manager.NewManagerService(topics)

	midF.AddMiddleware(apiToken, allowedOrings)
	router := router.New()

	// API
	fi.NewSystemAPI(router, system)
	fi.NewMonitorAPI(router, reg)
	fi.NewSwaggerAPI(router)
	fi.NewManagerAPI(router, manager)

	r := fasthttp.Server{
		Handler:                      midF.Middlewares(router.Handler),
		Name:                         "dominus-project",
		MaxRequestBodySize:           1000,
		ReduceMemoryUsage:            true,
		DisablePreParseMultipartForm: true,
		KeepHijackedConns:            true,
		CloseOnShutdown:              true,
		StreamRequestBody:            true,
		Logger:                       logs,
	}

	_, errC := os.Stat(env.SslCert)
	_, errK := os.Stat(env.KeyFile)
	_, errCa := os.Stat(env.SslCaCert)

	if errC == nil && errK == nil {
		go func(port int64, cert, key string, cancel context.CancelFunc) {
			log.Fatal(r.ListenAndServeTLS(fmt.Sprintf("0.0.0.0:%d", port), cert, key))
			cancel()
		}(env.RestPort, env.SslCert, env.KeyFile, cancel)
	} else {
		go func(port int64, cancel context.CancelFunc) {
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
		queue  = make(chan struct{})
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
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			metricserver.UnaryServerInterceptor(),
			auth.UnaryServerInterceptor(midGs.ApiToken),
			logging.UnaryServerInterceptor(midGs.LogErrors()),
		),
		grpc.ChainStreamInterceptor(
			metricserver.StreamServerInterceptor(),
			auth.StreamServerInterceptor(midGs.ApiToken),
			logging.StreamServerInterceptor(midGs.LogErrors()),
		),
	)

	optsD = append(optsD,
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithUnaryInterceptor(metricclient.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(metricclient.StreamClientInterceptor()),
		grpc.WithUnaryInterceptor(midGc.UnaryAuthInterceptor),
		grpc.WithStreamInterceptor(midGc.StreamAuthInterceptor),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`),
	)

	gclient := gt.NewGrpClient(optsD, logs)

	// Interactors
	grpcC := grpcconn.NewGrpcService(logs, gclient, topics)

	// API
	srG := gi.NewGrpcAPI(optsS, grpcC, queue)

	// Server
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", env.GrpcPort))
	if err != nil {
		log.Println(err)
		return
	}
	go func(sr *grpc.Server, list net.Listener, cancel context.CancelFunc) {
		log.Fatal(sr.Serve(list))
		cancel()
	}(srG, listener, cancel)

	//*********************************
	//**********Banner*****************
	//*********************************

	if errC == nil && errK == nil {
		if *showBanner {
			fmt.Printf("%s Rest: https://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", banner, env.RestPort, env.GrpcPort)
		} else {
			fmt.Printf("Rest: https://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", env.RestPort, env.GrpcPort)
		}
	} else {
		if *showBanner {
			fmt.Printf("%s Rest: http://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", banner, env.RestPort, env.GrpcPort)
		} else {
			fmt.Printf("Rest: http://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", env.RestPort, env.GrpcPort)
		}
	}

	//*********************************
	//*********Shutdown servers********
	//*********************************

	// Wait for a signal
	select {
	case err := <-st:
		log.Fatalln(err)
	case <-ctx.Done():
		log.Fatalln("[-] Context closed")
	}
	
	queue <- struct{}{}
	if err := r.Shutdown(); err != nil {
		log.Fatal(err)
	}
	srG.GracefulStop()
	close(queue)
	close(st)
}
