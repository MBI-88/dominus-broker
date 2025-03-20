package main

import (
	"context"
	"dominus-project/app/domain/config"
	"dominus-project/app/domain/rules"
	"dominus-project/app/interactors"
	"dominus-project/docs"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"

	"dominus-project/app/interfaces/events"
	fi "dominus-project/app/interfaces/fasthttp/input"
	fm "dominus-project/app/interfaces/fasthttp/middlewares"
	gi "dominus-project/app/interfaces/grpc/input"
	gm "dominus-project/app/interfaces/grpc/middlewares"
	gt "dominus-project/app/interfaces/grpc/output"
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
	system     chan os.Signal
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

dominus-project server is running on`
)

// catches inital variables
func init() {
	mode = flag.Bool("prod", false, "set operation mode")
	showBanner = flag.Bool("banner", true, "show banner")

	flag.Usage = func() {
		info := fmt.Sprintf("[*] ***dominus-project*** [*]\n")
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
// @version 1.0.0
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name API_TOKEN
// @host localhost:8000
// @BasePath /
func main() {
	//Receives commands from cli
	flag.Parse()

	//Instances
	settings := config.NewConfig()
	env := settings.GetEnvVar(*mode)
	rls := rules.NewRules()

	logs := events.NewLogs(env.Logs)
	inter := interactors.NewInteractor(logs, rls)
	docs.SwaggerInfo.Host = env.Host

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

	midF.AddMiddleware(apiToken, allowedOrings)
	router := router.New()
	fi.NewSystemAPI(router, inter.NewSystemService())
	fi.NewMonitorAPI(router, reg)
	fi.NewSwaggerAPI(router)

	r := fasthttp.Server{
		Handler:                            midF.Middlewares(router.Handler),
		Name:                               "dominus-project",
		MaxRequestBodySize:                 1000,
		ReduceMemoryUsage:                  true,
		DisablePreParseMultipartForm:       true,
		KeepHijackedConns:                  true,
		CloseOnShutdown:                    true,
		StreamRequestBody:                  true,
		Logger:                             logs,
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

		credsD, err := credentials.NewClientTLSFromFile(env.SslCaCert, "dominus-project.com")
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
			metricserver.UnaryServerInterceptor(),
			auth.UnaryServerInterceptor(midGs.ApiToken),
			logging.UnaryServerInterceptor(midGs.LogErrors()),
		),
		grpc.ChainStreamInterceptor(
			otelgrpc.StreamServerInterceptor(),
			metricserver.StreamServerInterceptor(),
			auth.StreamServerInterceptor(midGs.ApiToken),
			logging.StreamServerInterceptor(midGs.LogErrors()),
		),
	)

	optsD = append(optsD,
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithUnaryInterceptor(metricclient.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		grpc.WithStreamInterceptor(metricclient.StreamClientInterceptor()),
		grpc.WithUnaryInterceptor(midGc.UnaryAuthInterceptor),
		grpc.WithStreamInterceptor(midGc.StreamAuthInterceptor),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)),
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig":[{"round_robin":{}}]}`),
	)

	gclient := gt.NewGrpClient(optsD)
	inter = inter.Set(gclient)
	srG := gi.NewGrpcAPI(optsS, inter.NewGrpcService())
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
	case err := <-system:
		log.Fatalln(err)
	case <-ctx.Done():
		log.Fatalln("[-] Context closed")
	}

	if err := r.Shutdown(); err != nil {
		log.Fatal(err)
	}
	srG.GracefulStop()
}
