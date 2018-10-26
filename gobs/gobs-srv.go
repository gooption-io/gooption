package main

import (
	"context"
	"net"
	"net/http"

	"github.com/gooption-io/gooption/gobs/pb"
	"github.com/gooption-io/gooption/utils"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

type service struct {
	config     utils.ServiceConfig
	serverImpl pb.GobsServer
}

func NewService(impl pb.GobsServer, conf utils.ServiceConfig) *service {
	logrus.Infoln(conf)
	return &service{
		config:     conf,
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
	err := pb.RegisterGobsHandlerFromEndpoint(ctx, mux, s.config.TCP, opts)
	if err != nil {
		errCh <- err
	}

	logrus.Infoln("http server ready on port ", s.config.HTTP)
	errCh <- http.ListenAndServe(s.config.HTTP, cors.Default().Handler(mux))
}

// methods takes server chain as argument so it remains configurable per service while not changing core logic
// might be useful for dependency injection
func (s *service) ServeTCP(errCh chan error) {
	lis, err := net.Listen("tcp", s.config.TCP)
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

	logrus.Infoln("grpc server ready on port ", s.config.TCP)
	errCh <- grpcSrv.Serve(lis)
}

// methods takes server chain as argument so it remains configurable per service while not changing core logic
// might be useful for dependency injection
func (s *service) ServePromHTTP(errCh chan error) {
	logrus.Infoln("promhttp server ready on port ", s.config.PromHTTP)
	http.Handle("/metrics", promhttp.Handler())
	errCh <- http.ListenAndServe(s.config.PromHTTP, nil)
}

func (s *service) Serve() {
	errCh := make(chan error, 3)
	defer close(errCh)

	go s.ServeTCP(errCh)
	go s.ServeHTTP(errCh)
	go s.ServePromHTTP(errCh)

	for i := 0; i < 3; i++ {
		err := <-errCh
		if err != nil {
			logrus.Error(err)
		}
	}
}
