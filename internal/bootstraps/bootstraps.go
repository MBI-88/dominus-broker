package bootstraps

import (
	"context"
	"dominus-broker/config"

	"dominus-broker/internal/application/usecases/broker"
	"dominus-broker/internal/application/usecases/sqs"
	"dominus-broker/internal/infrastructure/enum"

	"dominus-broker/internal/infrastructure/event"
	fi "dominus-broker/internal/infrastructure/fasthttp/inbound"
	fm "dominus-broker/internal/infrastructure/fasthttp/middlewares"
	gi "dominus-broker/internal/infrastructure/grpc/inbound"
	gm "dominus-broker/internal/infrastructure/grpc/middlewares"
	gt "dominus-broker/internal/infrastructure/grpc/outbound"
	"dominus-broker/internal/infrastructure/redis/cchecker"
	"dominus-broker/internal/infrastructure/redis/cmemory"

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
	"google.golang.org/grpc/reflection"
)

func gRPServer(
	cf *config.Config,
	logs event.Event,
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

	checker := cchecker.NewCheckerClient(
		cf.RedisConfig.Port,
		cf.RedisConfig.CheckerDB,
		cf.RedisConfig.Host,
		cf.RedisConfig.Password,
		cf.RedisConfig.Tls,
		cf.RedisConfig.Username,
		cf.RedisConfig.IdPotencyEx,
	)

	midGs := gm.NewMiddleware(cf.GrpcConfig.ApiToken, logs, checker)
	midGc := gm.NewInterceptor(cf.GrpcConfig.ApiToken, logs)

	if errC == nil && errK == nil && errCa == nil {
		credsS, err := credentials.NewServerTLSFromFile(cf.CertConfig.SslCert, cf.CertConfig.KeyFile)
		if err != nil {
			cancel()
			panic(err)
		}
		optsS = append(optsS, grpc.Creds(credsS))

		credsD, err := credentials.NewClientTLSFromFile(cf.CertConfig.SslCaCert, enum.DOMAIN)
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
			auth.UnaryServerInterceptor(midGs.IdPotency),
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
		grpc.WithDefaultServiceConfig(enum.ROUTER_CLIENT),
	)

	// Clients
	bclient := gt.NewGrpClient(optsD, logs)
	qclient := cmemory.NewMemoryClient(
		cf.RedisConfig.Port,
		cf.RedisConfig.MemoryDB,
		cf.RedisConfig.Host,
		cf.RedisConfig.Password,
		cf.RedisConfig.Tls,
		cf.RedisConfig.Username,
		cf.RedisConfig.StreamID,
	)

	// Interactors
	broker := broker.NewBroker(bclient)
	sqs := sqs.NewSQS(qclient)

	// Create group if not exist
	if err := qclient.Group(cf.RedisConfig.GroupID); err != nil {
		logs.WriteLog(context.Background(), enum.ERROR, "qclient.Group", err.Error())
	}

	// Server
	server := grpc.NewServer(optsS...)

	// Register APIs
	gi.NewBrokerAPI(server, broker, logs)
	gi.NewSqsAPI(server, sqs, logs)

	// Reflecting server
	reflection.Register(server)

	// Listener
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cf.GrpcConfig.Port))
	if err != nil {
		log.Println(err)
		return nil
	}
	go func(sr *grpc.Server, list net.Listener, cancel context.CancelFunc) {
		log.Fatal(sr.Serve(list))
		cancel()
	}(server, listener, cancel)

	return server
}

func restServer(
	reg *prometheus.Registry,
	cf *config.Config,
	logs event.Event,
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

	r := fasthttp.Server{
		Handler:                      midF.Middlewares(router.Handler),
		Name:                         enum.PROJECT_NAME,
		MaxRequestBodySize:           100,
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
		}(cf.RestConfig.Port, cf.CertConfig.SslCert, cf.CertConfig.KeyFile, cancel)
	} else {
		go func(port int64, cancel context.CancelFunc) {
			log.Fatal(r.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port)))
			cancel()
		}(cf.RestConfig.Port, cancel)
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
	cf := config.NewConfig()

	lgs := event.NewEvent(cf.LogConfig.LogMode, cf.LogConfig.LogUrl, *mode)

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

	r := restServer(reg, cf, lgs, cancel, errC, errK)

	//********************************
	//*********Grpc server************
	//********************************

	srG := gRPServer(cf, lgs, errC, errK, errCa, cancel, metricserver, metricclient)
	if srG == nil {
		panic("GRPC server error")
	}

	//*********************************
	//**********Banner*****************
	//*********************************

	if errC == nil && errK == nil {
		if *showBanner {
			fmt.Printf("%s Rest: https://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", banner, cf.RestConfig.Port, cf.GrpcConfig.Port)
		} else {
			fmt.Printf("Rest: https://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", cf.RestConfig.Port, cf.GrpcConfig.Port)
		}
	} else {
		if *showBanner {
			fmt.Printf("%s Rest: http://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", banner, cf.RestConfig.Port, cf.GrpcConfig.Port)
		} else {
			fmt.Printf("Rest: http://0.0.0.0:%d 🚀  Grpc: 0.0.0.0:%d 🚀\n", cf.RestConfig.Port, cf.GrpcConfig.Port)
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
