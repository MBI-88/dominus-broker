package boostrap

import (
	"context"
	"dominus-project/config"
	"dominus-project/docs"
	"dominus-project/internal/domain/adapters"
	"dominus-project/internal/orchestrators/broker"
	"dominus-project/internal/orchestrators/queue"

	"dominus-project/internal/infrastructure/events"
	fi "dominus-project/internal/infrastructure/fasthttp/input"
	fm "dominus-project/internal/infrastructure/fasthttp/middlewares"
	gi "dominus-project/internal/infrastructure/grpc/input"
	gm "dominus-project/internal/infrastructure/grpc/middlewares"
	gt "dominus-project/internal/infrastructure/grpc/output"
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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
)

func gRPServer(
	cf config.Config,
	logs adapters.Logs,
	errC,
	errK,
	errCa error,
	cancel context.CancelFunc,
	metricserver *grpcmetrics.ServerMetrics,
	metricclient *grpcmetrics.ClientMetrics,
) *grpc.Server {
	var (
		optsS []grpc.ServerOption
		optsD []grpc.DialOption
	)

	midGs := gm.NewMiddleware(cf.GrpcConfig.ConnectionKey, logs)
	midGc := gm.NewInterceptor(cf.GrpcConfig.ConnectionKey, logs)

	if errC == nil && errK == nil && errCa == nil {
		credsS, err := credentials.NewServerTLSFromFile(cf.CertConfig.SslCert, cf.CertConfig.KeyFile)
		if err != nil {
			cancel()
			panic(err)
		}
		optsS = append(optsS, grpc.Creds(credsS))

		credsD, err := credentials.NewClientTLSFromFile(cf.CertConfig.SslCaCert, "dominus.com")
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
			auth.UnaryServerInterceptor(midGs.ApiToken),
			metricserver.UnaryServerInterceptor(),
			logging.UnaryServerInterceptor(midGs.LogErrors()),
		),
		grpc.ChainStreamInterceptor(
			auth.StreamServerInterceptor(midGs.ApiToken),
			metricserver.StreamServerInterceptor(),
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

	bclient := gt.NewGrpClient(optsD, logs)
	//qclient :=  add client implementation

	// Interactors
	brk := broker.NewBroker(logs, bclient)
	quk := queue.NewQueue(logs, nil, nil) 

	// API
	srGRP := gi.NewGrpcAPI(optsS, brk, quk, logs)

	// Server
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cf.GrpcConfig.GRPCPort))
	if err != nil {
		log.Println(err)
		return nil
	}
	go func(sr *grpc.Server, list net.Listener, cancel context.CancelFunc) {
		log.Fatal(sr.Serve(list))
		cancel()
	}(srGRP, listener, cancel)

	return srGRP
}

func restServer(
	reg *prometheus.Registry,
	cf config.Config,
	logs adapters.Logs,
	cancel context.CancelFunc,
	errC error,
	errK error,
) *fasthttp.Server {
	midF := fm.NewMiddleware()
	apiToken := fm.NewMiddlewareApiToken(cf.RestConfig.ApiToken, logs)
	allowedOrings := fm.NewMiddlewareHost(cf.RestConfig.AllowOrigins, logs)

	midF.AddMiddleware(apiToken, allowedOrings)
	router := router.New()

	// API
	fi.NewMonitorAPI(router, reg, logs)
	fi.NewSwaggerAPI(router)

	r := fasthttp.Server{
		Handler:                      midF.Middlewares(router.Handler),
		Name:                         "dominus-project",
		MaxRequestBodySize:           1000,
		ReduceMemoryUsage:            true,
		DisablePreParseMultipartForm: true,
		KeepHijackedConns:            true,
		CloseOnShutdown:              true,
		StreamRequestBody:            true,
	}

	if errC == nil && errK == nil {
		go func(port int64, cert, key string, cancel context.CancelFunc) {
			log.Fatal(r.ListenAndServeTLS(fmt.Sprintf("0.0.0.0:%d", port), cert, key))
			cancel()
		}(cf.RestConfig.RestPort, cf.CertConfig.SslCert, cf.CertConfig.KeyFile, cancel)
	} else {
		go func(port int64, cancel context.CancelFunc) {
			log.Fatal(r.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port)))
			cancel()
		}(cf.RestConfig.RestPort, cancel)
	}

	return &r
}

func monitoring() (*prometheus.Registry, *grpcmetrics.ServerMetrics, *grpcmetrics.ClientMetrics) {
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

	return reg, metricserver, metricclient
}

func RunApp(mode, showBanner *bool, banner string) {
	//Instances
	cf := config.NewConfig(*mode)

	logs := events.NewLogs(cf.InfraConfig.LogMode, cf.InfraConfig.LogURL, *mode)
	docs.SwaggerInfo.Host = cf.InfraConfig.Host

	// Signals
	st := make(chan os.Signal, 1)
	signal.Notify(st, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(st)

	//context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, errC := os.Stat(cf.CertConfig.SslCert)
	_, errK := os.Stat(cf.CertConfig.KeyFile)
	_, errCa := os.Stat(cf.CertConfig.SslCaCert)

	//**************************************
	//***********Monitoring*****************
	//**************************************

	reg, metricserver, metricclient := monitoring()

	//**************************************
	//***********Rest Server****************
	//**************************************

	r := restServer(reg, cf, logs, cancel, errC, errK)

	//********************************
	//*********Grpc server************
	//********************************

	srG := gRPServer(cf, logs, errC, errK, errCa, cancel, metricserver, metricclient)
	if srG == nil {
		panic("GRPC server error")
	}

	//*********************************
	//**********Banner*****************
	//*********************************

	if errC == nil && errK == nil {
		if *showBanner {
			fmt.Printf("%s Rest: https://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", banner, cf.RestConfig.RestPort, cf.GrpcConfig.GRPCPort)
		} else {
			fmt.Printf("Rest: https://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", cf.RestConfig.RestPort, cf.GrpcConfig.GRPCPort)
		}
	} else {
		if *showBanner {
			fmt.Printf("%s Rest: http://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", banner, cf.RestConfig.RestPort, cf.GrpcConfig.GRPCPort)
		} else {
			fmt.Printf("Rest: http://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", cf.RestConfig.RestPort, cf.GrpcConfig.GRPCPort)
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

	if err := r.Shutdown(); err != nil {
		log.Fatal(err)
	}
	srG.GracefulStop()
	close(st)
}
