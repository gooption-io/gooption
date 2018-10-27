package main

import (
	"flag"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gooption-io/gooption/proto/go/pb"
	"github.com/gooption-io/gooption/utils"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	context "golang.org/x/net/context"
	"gonum.org/v1/gonum/stat/distuv"
)

const (
	ivSeed    = 0.1 // solver starting point
	maxIter   = 1000
	putLBound = 0.20
)

var (
	phi  = distuv.Normal{Mu: 0, Sigma: 1}.CDF
	dphi = distuv.Normal{Mu: 0, Sigma: 1}.Prob

	allGreeks  = []string{"delta", "gamma", "vega", "theta", "rho"}
	putCallMap = map[string]float64{
		"call": 1.0,
		"put":  -1.0,
	}

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

	// // prom
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

	var (
		mult = putCallMap[strings.ToLower(in.Contract.Putcall)]

		s = in.Marketdata.Spot.Index.Value
		v = in.Marketdata.Vol.Index.Value
		r = in.Marketdata.Rate.Index.Value
		k = in.Contract.Strike
		t = time.Unix(int64(in.Contract.Expiry), 0).Sub(
			time.Unix(int64(in.Pricingdate), 0)).Hours() / 24.0 / 365.250
	)

	return &pb.PriceResponse{
		Price: bs(s, v, r, k, t, mult),
	}, nil
}

/*
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
*/
func bs(s, v, r, k, t, mult float64) float64 {
	d1 := d1(s, k, t, v, r)
	d2 := d2(d1, v, t)

	return mult * (s*phi(mult*d1) - k*phi(mult*d2)*math.Exp(-r*t))
}

func d1(S, K, T, Sigma, R float64) float64 {
	return (1.0 / Sigma * math.Sqrt(T)) * (math.Log(S/K) + (R+Sigma*Sigma*0.5)*T)
}

func d2(d1, Sigma, T float64) float64 {
	return d1 - Sigma*math.Sqrt(T)
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

	var (
		mult = putCallMap[strings.ToLower(in.Request.Contract.Putcall)]

		s = in.Request.Marketdata.Spot.Index.Value
		v = in.Request.Marketdata.Vol.Index.Value
		r = in.Request.Marketdata.Rate.Index.Value
		k = in.Request.Contract.Strike
		t = time.Unix(int64(in.Request.Contract.Expiry), 0).Sub(
			time.Unix(int64(in.Request.Pricingdate), 0)).Hours() / 24.0 / 365.250

		d1 = d1(s, k, t, v, r)
		d2 = d2(d1, v, t)
	)

	greeks := make([]*pb.GreekResponse_Greek, len(in.Greek))
	for index, greek := range in.Greek {
		greekResponse := &pb.GreekResponse_Greek{
			Label: greek,
		}

		switch greek {
		case "delta":
			greekResponse.Value = delta(d1, mult)
		case "gamma":
			greekResponse.Value = gamma(s, t, v, d1)
		case "vega":
			greekResponse.Value = vega(s, t, d1)
		case "theta":
			greekResponse.Value = theta(s, k, t, v, r, d1, d2, mult)
		case "rho":
			greekResponse.Value = rho(k, t, r, d2, mult)
		default:
			greekResponse.Error = "Unknown greek " + greek
		}

		greeks[index] = greekResponse
	}

	return &pb.GreekResponse{
		Greeks: greeks,
	}, nil
}

func delta(d1, mult float64) float64 {
	return mult * phi(mult*d1)
}

func gamma(s, t, sigma, d1 float64) float64 {
	return dphi(d1) / (s * sigma * math.Sqrt(t))
}

func vega(s, t, d1 float64) float64 {
	return s * dphi(d1) * math.Sqrt(t)
}

func theta(s, k, t, sigma, r, d1, d2, mult float64) float64 {
	return -0.5*(s*dphi(d1)*sigma/math.Sqrt(t)) - (mult * r * k * math.Exp(-r*t) * phi(mult*d2))
}

func rho(k, t, r, d2, mult float64) float64 {
	return mult * k * t * math.Exp(-r*t) * phi(mult*d2)
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

type ivSolverResult struct {
	IV                float64
	NbSolverIteration int
}

func bsImpliedVol(index int, in *pb.ImpliedVolRequest, out chan<- pb.ImpliedVolSlice) {
	var (
		slice = in.Quotes[index]
		mult  = func(q *pb.OptionQuote) float64 { return putCallMap[strings.ToLower(q.Putcall)] }

		s = in.Marketdata.Spot.Index.Value
		r = in.Marketdata.Rate.Index.Value
		k = func(q *pb.OptionQuote) float64 { return q.Strike }
		p = func(q *pb.OptionQuote) float64 { return (q.Ask + q.Bid) / 2.0 }
		t = func(q *pb.OptionQuoteSlice) float64 {
			return time.Unix(int64(q.Expiry), 0).Sub(time.Unix(int64(in.Pricingdate), 0)).Hours() / 24.0 / 365.250
		}
	)

	calibratedSlice := pb.ImpliedVolSlice{
		Timestamp: in.Marketdata.Timestamp,
		Expiry:    slice.Expiry,
		Iserror:   false,
		Quotes:    make([]*pb.ImpliedVolQuote, len(slice.Puts)+len(slice.Calls)),
	}

	calibrate := func(quote *pb.OptionQuote) *pb.ImpliedVolQuote {
		res, err := ivRootSolver(p(quote), s, r, k(quote), t(slice), mult(quote))
		if err != nil {
			calibratedSlice.Iserror = true
			return &pb.ImpliedVolQuote{
				Timestamp:   quote.Timestamp,
				Input:       quote,
				Error:       err.Error(),
				Nbiteration: int64(res.NbSolverIteration),
			}
		}

		return &pb.ImpliedVolQuote{
			Timestamp:   quote.Timestamp,
			Input:       quote,
			Vol:         res.IV,
			Nbiteration: int64(res.NbSolverIteration),
		}
	}

	calibIndex := 0
	for _, put := range slice.Puts {
		if put.Strike/s > putLBound && put.Strike/s <= 1.0 {
			calibratedSlice.Quotes[calibIndex] = calibrate(put)
			calibIndex++
		}
	}

	for _, call := range slice.Calls {
		if call.Strike/s > 1.0 {
			calibratedSlice.Quotes[calibIndex] = calibrate(call)
			calibIndex++
		}
	}

	calibratedSlice.Quotes = calibratedSlice.Quotes[0:calibIndex]
	out <- calibratedSlice
}

/*
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
*/
func ivRootSolver(mktPrice, s, r, k, t, mult float64) (*ivSolverResult, error) {
	var (
		iv      = 0.1
		maxIter = 1000
	)

	for index := 0; index < maxIter; index++ {
		bsPrice := bs(s, iv, r, k, t, mult)
		iv = iv - (bsPrice-mktPrice)/vega(s, t, d1(s, k, t, iv, r))
		if math.Abs(bsPrice-mktPrice) < 1E-10 { //decrease to 1E-25 to test convergence error
			return &ivSolverResult{iv, index}, nil
		}
	}
	return &ivSolverResult{iv, maxIter}, errors.New("Did not converge to required interval")
}
