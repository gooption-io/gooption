package main

import (
	"context"
	"net"
	"net/http"

	"github.com/gooption-io/gooption/gobs/pb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

var tcpPort, httpPort, promhttpPort string
var tcpReqs *prometheus.CounterVec

func init() {
	tcpPort = *flag.String("tcp-listen-address", ":50051", "The Port to listen on for TCP requests")
	httpPort = *flag.String("http-listen-address", ":8081", "The Port to listen on for HTTP requests")
	promhttpPort = *flag.String("prom-listen-address", ":8080", "The Port to listen on for Promhttp requests")
	tcpReqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tcp_requests_total",
			Help: "How many TCP requests processed, partitioned by request type",
		},
		[]string{"code"},
	)

	prometheus.MustRegister(tcpReqs)
}

type service struct {
	serverImpl pb.GobsServer
}

func NewService(impl pb.GobsServer) *service {
	return &service{
		serverImpl: impl,
	}
}

// methods takes server chain as argument so it remains configurable per service while not changing core logic
// might be useful for dependency injection
func (s *service) ServeHTTP(errCh chan error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterGobsHandlerFromEndpoint(ctx, mux, tcpPort, opts)
	if err != nil {
		errCh <- err
	}

	logrus.Infoln("http server ready on port ", httpPort)
	errCh <- http.ListenAndServe(httpPort, cors.Default().Handler(mux))
}

// methods takes server chain as argument so it remains configurable per service while not changing core logic
// might be useful for dependency injection
func (s *service) ServeTCP(errCh chan error) {
	lis, err := net.Listen("tcp", tcpPort)
	if err != nil {
		errCh <- err
	}
	defer lis.Close()

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(func(methodFullName string, err error) bool {
			// will not log gRPC calls if it was a call to healthcheck and no error was raised
			if err == nil && methodFullName == "main.server.healthcheck" {
				return false
			}

			// by default you will log all calls
			return true
		}),
	}

	grpcSrv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New()), opts...)))
	pb.RegisterGobsServer(grpcSrv, s.serverImpl)
	reflection.Register(grpcSrv)

	logrus.Infoln("grpc server ready on port ", tcpPort)
	errCh <- grpcSrv.Serve(lis)
}

// methods takes server chain as argument so it remains configurable per service while not changing core logic
// might be useful for dependency injection
func (s *service) ServePromHTTP(errCh chan error) {
	logrus.Infoln("promhttp server ready on port ", promhttpPort)
	http.Handle("/metrics", promhttp.Handler())
	errCh <- http.ListenAndServe(promhttpPort, nil)
}

func (s *service) Serve() {
	errCh := make(chan error, 3)
	defer close(errCh)

	go s.ServeTCP(errCh)
	go s.ServeTCP(errCh)
	go s.ServePromHTTP(errCh)

	for i := 0; i < 3; i++ {
		err := <-errCh
		if err != nil {
			logrus.Error(err)
		}
	}
}
