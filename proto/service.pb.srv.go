package gooption

import (
	"log"
	"net"
	"net/http"

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

func ServeEuropeanOptionPricerServer(tcp, prom string, server EuropeanOptionPricerServer) error {
	lis, err := net.Listen("tcp", tcp)
	if err != nil {
		return err
	}
	defer lis.Close()

	grpcSrv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			grpc.UnaryServerInterceptor(grpc_prometheus.UnaryServerInterceptor),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New()))))

	RegisterEuropeanOptionPricerServer(grpcSrv, server)

	reflection.Register(grpcSrv)
	grpc_prometheus.Register(grpcSrv)

	httpServer := &http.Server{Handler: promhttp.Handler(), Addr: prom}
	go func() {
		// Start your http server for prometheus.
		logrus.Infoln("EuropeanOptionPricer prometheus server ready at ", prom)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	logrus.Infoln("EuropeanOptionPricer grpc server ready on port at ", tcp)
	return grpcSrv.Serve(lis)
}

func ServeEuropeanOptionPricerServerGateway(tcpPort, httpPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := RegisterEuropeanOptionPricerHandlerFromEndpoint(ctx, mux, tcpPort, opts)
	if err != nil {
		return err
	}

	logrus.Infoln("connected to EuropeanOptionPricer at ", tcpPort)
	logrus.Infoln("EuropeanOptionPricer reverse proxy ready at ", httpPort)
	return http.ListenAndServe(httpPort, cors.Default().Handler(mux))
}
