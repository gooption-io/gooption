//go:generate sh -c "protoc --proto_path=pb --proto_path=$GOPATH/src/github.com/lehajam/gooption/pb --proto_path=$GOPATH/src --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogofast_out=plugins=grpc:pb $GOPATH/src/github.com/lehajam/gooption/pb/*.proto pb/*.proto"
//go:generate sh -c "protoc --proto_path=pb --proto_path=$GOPATH/src/github.com/lehajam/gooption/pb --proto_path=$GOPATH/src --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:pb pb/*.proto"
// -- //go:generate sh -c "protoc --proto_path=pb --proto_path=$GOPATH/src/github.com/lehajam/gooption/pb --proto_path=$GOPATH/src --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogofast_out=plugins=grpc:pb --grpc-gateway_out=logtostderr=true:pb $GOPATH/src/github.com/lehajam/gooption/pb/*.proto pb/*.proto"
//go:generate gooption-cli -p gobs -r Price -r Greek -r ImpliedVol
package main

import (
	"net/http"
	"strings"
	"sync"

	"github.com/namsral/flag"

	context "golang.org/x/net/context"

	"errors"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/lehajam/gooption/goquantlib/pb"
	"github.com/lehajam/gooption/goquantlib/quantlib"
)

var (
	config       = flag.String(flag.DefaultConfigFlagname, "", "The Path to config file")
	tcpPort      = flag.String("tcp-listen-address", ":50051", "The Port to listen on for TCP requests")
	httpPort     = flag.String("http-listen-address", ":8081", "The Port to listen on for HTTP requests")
	promhttpPort = flag.String("prom-listen-address", ":8080", "The Port to listen on for Promhttp requests")

	putCallMap = map[string]float64{
		"call": 1.0,
		"put":  -1.0,
	}

	httpReqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)

	tcpReqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tcp_requests_total",
			Help: "How many TCP requests processed, partitioned by return status",
		},
		[]string{"code"},
	)
)

// server is used to implement pb.ModerlServer.
type server struct{}

// methods takes server chain as argument so it remains configurable per service while not changing core logic
// might be useful for dependency injection
func httpServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterGobsHandlerFromEndpoint(ctx, mux, *tcpPort, opts)
	if err != nil {
		return err
	}

	logrus.Infoln("http server ready on port ", *httpPort)
	return http.ListenAndServe(*httpPort, cors.Default().Handler(mux))
}

// methods takes server chain as argument so it remains configurable per service while not changing core logic
// might be useful for dependency injection
func tcpServer() error {
	lis, err := net.Listen("tcp", *tcpPort)
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

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New()), opts...)))
	pb.RegisterGobsServer(s, &server{})
	reflection.Register(s)

	logrus.Infoln("grpc server ready on port ", *tcpPort)
	return s.Serve(lis)
}

func promhttpServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterGobsHandlerFromEndpoint(ctx, mux, *tcpPort, opts)
	if err != nil {
		return err
	}

	logrus.Infoln("promhttp server ready on port ", *promhttpPort)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(*promhttpPort, nil)
}

func start(entrypoint func() error) {
	defer glog.Flush()
	if err := entrypoint(); err != nil {
		logrus.Fatal(err)
		panic(err)
	}
}

func init() {
	prometheus.MustRegister(httpReqs)
	prometheus.MustRegister(tcpReqs)
}

func ExampleCounterVec() {
	// A implementer dans l'appel du service HTTP
	httpReqs.WithLabelValues("404", "POST").Add(54)
	httpReqs.WithLabelValues("200", "POST").Add(450003)
	// A implémenter dans l'appel du servive GRPC
	tcpReqs.WithLabelValues("OK").Add(150)
	tcpReqs.WithLabelValues("KO").Add(12)
}

func main() {
	flag.Parse()

	// A supprimer une fois les compteurs implementés
	ExampleCounterVec()

	var wg sync.WaitGroup
	wg.Add(3)
	go start(tcpServer)
	go start(httpServer)
	go start(promhttpServer)
	wg.Wait()
}

/*
Price computes the fair value of a european stock option according to Black Scholes formula
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
*/
func (srv *server) Price(ctx context.Context, in *pb.PriceRequest) (*pb.PriceResponse, error) {
	var (
		mult = int(putCallMap[strings.ToLower(in.Contract.Putcall)])
		s    = in.Marketdata.Spot.Index.Value
		v    = in.Marketdata.Vol.Index.Value
		r    = in.Marketdata.Rate.Index.Value
		k    = in.Contract.Strike
		p    = quantlib.EuropeanFlatVol(
			s,
			r,
			0,
			v,
			k,
			int(in.Pricingdate),
			int(in.Contract.Expiry),
			mult)
	)

	return &pb.PriceResponse{
		Price: p,
	}, nil
}

/*
Greeks computes the greeks of a european option according to Black Scholes formula
Black Scholes Greeks : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#The_Greeks
Possible values for Requests :  "all", "delta", "gamma", "vega", "theta", "rho"
Setting Request to "all" will compute all greeks
*/
func (srv *server) Greek(ctx context.Context, in *pb.GreekRequest) (*pb.GreekResponse, error) {
	logrus.Debugf("%+v\n", proto.MarshalTextString(in))
	return &pb.GreekResponse{}, errors.New("Not Implemented")
}

/*
ImpliedVol computes volatility matching the option quote passed as Quote using Newton Raphson solver
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
The second argument returned is the number of iteration used to converge
*/
func (srv *server) ImpliedVol(ctx context.Context, in *pb.ImpliedVolRequest) (*pb.ImpliedVolResponse, error) {
	logrus.Debugf("%+v\n", proto.MarshalTextString(in))
	return &pb.ImpliedVolResponse{}, errors.New("Not Implemented")
}
