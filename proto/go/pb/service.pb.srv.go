package pb

import (
	//"log"
	"net"
	"net/http"

	"github.com/gooption-io/gooption/v1/logging"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
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
)

func ServeEuropeanOptionPricerServer(tcp, prom string, server EuropeanOptionPricerServer) error {
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
			grpc_prometheus.UnaryServerInterceptor,
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New()), opts...)))

	RegisterEuropeanOptionPricerServer(grpcSrv, server)
	reflection.Register(grpcSrv)

	// prom
	reg := prometheus.NewRegistry()
	grpc_prometheus.Register(grpcSrv)
	reg.MustRegister(grpc_prometheus.NewServerMetrics())

	// Entrypoint HTTP handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, errWelcomePage := w.Write([]byte(`<html>
             <head><title>Go Option GOBS Prom</title></head>
             <body>
             <h1>Go Option GOBS Prometheus Metrics</h1>
             <p><a href="/metrics">Metrics</a></p>
             </body>
             </html>`))

		if errWelcomePage != nil {
			logging.Log(
				"error",
				"Failed to render Prometheus metrics welcome page %v", err)
		}
	})

	// HTTP Handler for Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		// Start your http server for prometheus.
		logging.Log("info", "EuropeanOptionPricer prometheus server ready at port %v", prom)
		if err := http.ListenAndServe(prom, nil); err != nil {
			logging.Log("fatal", "Unable to start Prometheus HTTP metrics server with error %v.", err)
		}
	}()

	logging.Log("info", "EuropeanOptionPricer grpc server ready on port at %v", tcp)
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

	logging.Log("info", "EuropeanOptionPricer reverse proxy ready at %v", httpPort)
	return http.ListenAndServe(httpPort, cors.Default().Handler(mux))
}
