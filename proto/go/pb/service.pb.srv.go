package pb

import (
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
)

func RunGobsServer(tcp, prom string, server GobsServer) error {
	lis, err := net.Listen("tcp", tcp)
	if err != nil {
		return err
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
			grpc_recovery.UnaryServerInterceptor(),
			grpc.UnaryServerInterceptor(grpc_prometheus.UnaryServerInterceptor),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New()), opts...)))

	RegisterGobsServer(grpcSrv, server)
	reflection.Register(grpcSrv)

	// prom
	reg := prometheus.NewRegistry()
	grpc_prometheus.Register(grpcSrv)
	reg.MustRegister(grpc_prometheus.NewServerMetrics())
	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: prom}
	go func() {
		// Start your http server for prometheus.
		logrus.Infoln("prometheus server ready on port ", prom)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	logrus.Infoln("grpc server ready on port ", tcp)
	return grpcSrv.Serve(lis)
}

func RunGobsGateway(tcpPort, httpPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := RegisterGobsHandlerFromEndpoint(ctx, mux, tcpPort, opts)
	if err != nil {
		return err
	}

	logrus.Infoln("http server ready on port ", httpPort)
	return http.ListenAndServe(httpPort, cors.Default().Handler(mux))
}
