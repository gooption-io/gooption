// go:generate sh -c "protoc --proto_path=$GOPATH/src/github.com/gooption/pb --proto_path=$GOPATH/src/github.com/gooption/gobs/pb --proto_path=$GOPATH/src/github.com/gooption/gobs/pb --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gofast_out=plugins=grpc:pb $GOPATH/src/github.com/gooption/gobs/pb/service.proto $GOPATH/src/github.com/gooption/pb/*.proto"
// go:generate sh -c "protoc --proto_path=$GOPATH/src/github.com/gooption/pb --proto_path=$GOPATH/src/github.com/gooption/gobs/pb --proto_path=$GOPATH/src/github.com/gooption/gobs/pb --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:pb $GOPATH/src/github.com/gooption/gobs/pb/service.proto"
// go:generate gooption-cli -p gobs -r Price -r Greek -r ImpliedVol
package main

import (
	"flag"
	"net/http"
	"sync"

	context "golang.org/x/net/context"

	"errors"
	"net"
	"sort"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/golang/glog"
	"github.com/gooption/gobs/pb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var (
	tcpPort      = ":50051"
	httpPort     = ":8081"
	gobsEndpoint = flag.String(
		"gobs_endpoint",
		"localhost:50051",
		"endpoint of YourService")
)

// server is used to implement pb.ModerlServer.
type server struct{}

func servehttp() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterGoBSServerHandlerFromEndpoint(ctx, mux, *gobsEndpoint, opts)
	if err != nil {
		return err
	}

	glog.V(2).Infoln("server ready on port %%s", httpPort)
	return http.ListenAndServe(httpPort, mux)
}

func servetcp() error {
	lis, err := net.Listen("tcp", tcpPort)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterGoBSServerServer(s, &server{})
	reflection.Register(s)

	glog.V(2).Infoln("server ready on port %s", tcpPort)
	return s.Serve(lis)
}

func start(entrypoint func() error) {
	defer glog.Flush()
	if err := entrypoint(); err != nil {
		glog.Fatal(err)
	}
}

func main() {
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(2)
	go start(servetcp)
	go start(servehttp)
	wg.Wait()
}

/*
Price computes the fair value of a european stock option according to Black Scholes formula
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
*/
func (srv *server) Price(ctx context.Context, in *pb.PriceRequest) (*pb.PriceResponse, error) {
	var (
		mult = putCallMap[in.Contract.Putcall]

		s = in.Marketdata.Spot.Index.Value
		v = in.Marketdata.Vol.Index.Value
		r = in.Marketdata.Rate.Index.Value
		k = in.Contract.Strike
		t = time.Unix(int64(in.Contract.Expiry), 0).Sub(
			time.Unix(int64(in.Pricingdate), 0)).Hours() / 24.0 / 365.250
		bs = bs(s, v, r, k, t, mult)
	)

	return &pb.PriceResponse{
		Price: bs,
	}, nil
}

/*
Greeks computes the greeks of a european option according to Black Scholes formula
Black Scholes Greeks : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#The_Greeks
Possible values for Requests :  "all", "delta", "gamma", "vega", "theta", "rho"
Setting Request to "all" will compute all greeks
*/
func (srv *server) Greek(ctx context.Context, in *pb.GreekRequest) (*pb.GreekResponse, error) {
	var (
		mult = putCallMap[in.Request.Contract.Putcall]

		s = in.Request.Marketdata.Spot.Index.Value
		v = in.Request.Marketdata.Vol.Index.Value
		r = in.Request.Marketdata.Rate.Index.Value
		k = in.Request.Contract.Strike
		t = time.Unix(int64(in.Request.Contract.Expiry), 0).Sub(
			time.Unix(int64(in.Request.Pricingdate), 0)).Hours() / 24.0 / 365.250
		d1 = d1(s, k, t, v, r)
		d2 = d2(d1, v, t)
	)

	if len(in.Greek) == 0 {
		return nil, errors.New("No greeks requested")
	}

	sort.Strings(in.Greek)
	if sort.SearchStrings(in.Greek, "all") < len(in.Greek) {
		in.Greek = allGreeks
	}

	response := &pb.GreekResponse{
		Greeks: make([]*pb.GreekResponse_Greek, len(in.Greek)),
	}

	for index := 0; index < len(in.Greek); index++ {
		response.Greeks[index] = &pb.GreekResponse_Greek{
			Label: in.Greek[index],
		}

		greek, err := bsGreek(response.Greeks[index].Label, s, v, r, k, t, mult, d1, d2)
		if err != nil {
			response.Greeks[index].Error = err.Error()
		}
		response.Greeks[index].Value = greek
	}

	return response, nil
}

/*
ImpliedVol computes volatility matching the option quote passed as Quote using Newton Raphson solver
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
The second argument returned is the number of iteration used to converge
*/
func (srv *server) ImpliedVol(ctx context.Context, in *pb.ImpliedVolRequest) (*pb.ImpliedVolResponse, error) {
	var (
		mult = func(q *pb.OptionQuote) float64 { return putCallMap[q.Putcall] }

		s = in.Marketdata.Spot.Index.Value
		r = in.Marketdata.Rate.Index.Value
		k = func(q *pb.OptionQuote) float64 { return q.Strike }
		p = func(q *pb.OptionQuote) float64 { return (q.Ask + q.Bid) / 2.0 }
		t = func(q *pb.OptionQuoteSlice) float64 {
			return time.Unix(int64(q.Expiry), 0).Sub(time.Unix(int64(in.Pricingdate), 0)).Hours() / 24.0 / 365.250
		}
	)

	surf := &pb.ImpliedVolSurface{
		Timestamp:  in.Marketdata.Timestamp,
		Volsurface: make([]*pb.ImpliedVolSlice, len(in.Quotes)),
	}

	for idx := 0; idx < len(in.Quotes); idx++ {
		go func(index int) {
			newOptionQuoteSliceIterator(in.Quotes[index], in.Marketdata).foreach(
				func(quote *pb.OptionQuote) *ivSolverResult {
					return ivRootSolver(p(quote), s, r, k(quote), t(in.Quotes[index]), mult(quote))
				}).then(
				func(calibratedSlice *pb.ImpliedVolSlice) {
					surf.Volsurface[index] = calibratedSlice
				})
		}(idx)
	}

	return &pb.ImpliedVolResponse{
		Volsurface: surf,
	}, nil
}

// mutation {
//   set {
//     _:stockidx <timestamp> "1.5066922e+09" .
//     _:stockidx <ticker> "AAPL" .
//     _:stockidx <value> "100.0" .
//     _:spot <index> _:stockidx .

//     _:volidx <timestamp> "1.5066922e+09" .
//     _:volidx <ticker> "AAPL" .
//     _:volidx <value> "0.10" .
//     _:vol <index> _:volidx .

//     _:ois <timestamp> "1.5066922e+09" .
//     _:ois <ticker> "USD.FEDFUND" .
//     _:ois <value> "0.01" .
//     _:rate <index> _:ois .

//     _:option <ticker> "AAPL DEC2017 PUT" .
// 		_:option <strike> "100.0" .
// 		_:option <expiry> "1.5092878e+09" .
//     _:option <putcall> "1" .

// 	   _:mkt <timestamp> "1.5066922e+09" .
//     _:mkt <spot> _:spot .
//     _:mkt <vol> _:vol .
//     _:mkt <rate> _:rate .

//     _:req <pricingdate> "1.5066922e+09" .
//     _:req <contract> _:option .
//     _:req <marketdata> _:mkt .
//   }
// }
//
// {
// 	"data": {
// 	  "code": "Success",
// 	  "message": "Done",
// 	  "uids": {
// 		"mkt": "0x7",
// 		"ois": "0x4",
// 		"option": "0x6",
// 		"rate": "0x5",
// 		"req": "0xa",
// 		"spot": "0x9",
// 		"stockidx": "0x8",
// 		"vol": "0x3",
// 		"volidx": "0x2"
// 	  }
// 	},
// 	"extensions": {
// 	  "server_latency": {
// 		"json": "57ms",
// 		"parsing": "252µs",
// 		"processing": "0s",
// 		"total": "57ms"
// 	  }
// 	}
//   }

// mutation{
// 	schema{
// 		  ticker: string @index(fulltext) .
// 	}

// 	set {
// 	  _:contract <ticker> "AAPL DEC2017 CALL 1200" .
// 		  _:contract <strike> "100" .
// 		  _:contract <expiry> "1509287800" .
// 	}
//   }
