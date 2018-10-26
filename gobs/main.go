//go:generate sh -c "protoc --proto_path=pb --proto_path=$GOPATH/src/github.com/lehajam/gooption/pb --proto_path=$GOPATH/src --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogofast_out=plugins=grpc:pb --grpc-gateway_out=logtostderr=true:pb $GOPATH/src/github.com/lehajam/gooption/pb/*.proto pb/*.proto"
//go:generate gooption-cli -p gobs -r Price -r Greek -r ImpliedVol
package main

import (
	"flag"

	"github.com/prometheus/client_golang/prometheus"

	context "golang.org/x/net/context"

	"sort"

	"github.com/pkg/errors"

	"github.com/gogo/protobuf/proto"

	"github.com/gooption-io/gooption/gobs/pb"
	"github.com/gooption-io/gooption/utils"
	"github.com/sirupsen/logrus"
)

var (
	env     = flag.String("env", "prod", "dev/prod config for ports")
	tcpReqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tcp_requests_total",
			Help: "How many TCP requests processed, partitioned by request type",
		},
		[]string{"code"},
	)
)

// server panic recovery
// for example if we fail to load config
func recoverServer() {
	if r := recover(); r != nil {
		logrus.WithField("error", r).Errorln("panic recovered")
	}
}

func main() {
	// recovery
	defer recoverServer()

	// flag
	flag.Parse()

	// prom
	prometheus.MustRegister(tcpReqs)

	// config
	utils.InitViperConfig("gobs", ".")

	// serve
	NewService(&server{}, utils.NewServiceConfig(*env)).Serve()
}

// server is used to implement pb.ModerlServer.
type server struct{}

/*
Price computes the fair value of a european stock option according to Black Scholes formula
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
*/
func (srv *server) Price(ctx context.Context, in *pb.PriceRequest) (*pb.PriceResponse, error) {
	tcpReqs.WithLabelValues("PriceRequest").Add(1)

	return &pb.PriceResponse{
		Price: bs(in.Pricingdate, in.Contract, in.Marketdata),
	}, nil
}

/*
Greeks computes the greeks of a european option according to Black Scholes formula
Black Scholes Greeks : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#The_Greeks
Possible values for Requests :  "all", "delta", "gamma", "vega", "theta", "rho"
Setting Request to "all" will compute all greeks
*/
func (srv *server) Greek(ctx context.Context, in *pb.GreekRequest) (*pb.GreekResponse, error) {
	tcpReqs.WithLabelValues("GreekRequest").Add(1)

	if len(in.Greek) == 0 {
		return nil, errors.New("No greeks requested")
	}

	sort.Strings(in.Greek)
	if sort.SearchStrings(in.Greek, "all") < len(in.Greek) {
		in.Greek = allGreeks
	}

	return &pb.GreekResponse{
		Greeks: bsGreek(in),
	}, nil
}

/*
ImpliedVol computes volatility matching the option quote passed as Quote using Newton Raphson solver
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
The second argument returned is the number of iteration used to converge
*/
func (srv *server) ImpliedVol(ctx context.Context, in *pb.ImpliedVolRequest) (*pb.ImpliedVolResponse, error) {
	tcpReqs.WithLabelValues("ImpliedVolRequest").Add(1)
	logrus.Debugf("%+v\n", proto.MarshalTextString(in))

	surf := &pb.ImpliedVolSurface{
		Timestamp: in.Marketdata.Timestamp,
		Slices:    make([]*pb.ImpliedVolSlice, len(in.Quotes)),
	}

	out := make(chan pb.ImpliedVolSlice, len(in.Quotes))
	defer close(out)

	for idx := 0; idx < len(in.Quotes); idx++ {
		go bsImpliedVol(idx, in, out)
	}

	for idx := 0; idx < len(in.Quotes); idx++ {
		slice := <-out
		surf.Slices[idx] = &slice
	}

	logrus.Debugf("%+v\n", proto.MarshalTextString(surf))
	return &pb.ImpliedVolResponse{
		Volsurface: surf,
	}, nil
}
