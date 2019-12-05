package server

import (
	"context"
	"gonum.org/v1/gonum/stat/distuv"
	"math"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	api_pb "gobs/api"
)

const (
	ivSeed    = 0.1 // solver starting point
	maxIter   = 1000
	putLBound = 0.20
)

var (

	phi        = distuv.Normal{Mu: 0, Sigma: 1}.CDF
	dphi       = distuv.Normal{Mu: 0, Sigma: 1}.Prob
	putCallMap = map[string]float64{"call": 1.0, "put": -1.0}
	allGreeks  = []string{"delta", "gamma", "vega", "theta", "rho"}
)

// EuropeanPricerServiceServer is a composite interface of api_pb.EuropeanPricerServiceServer and grapiserver.Server.
type EuropeanPricerServiceServer interface {
	api_pb.EuropeanPricerServiceServer
	grapiserver.Server
}

// NewEuropeanPricerServiceServer creates a new EuropeanPricerServiceServer instance.
func NewEuropeanPricerServiceServer() EuropeanPricerServiceServer {
	return &europeanPricerServiceServerImpl{}
}

type europeanPricerServiceServerImpl struct {
}

/*
Price computes the fair value of a european stock option according to Black Scholes formula
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends

Greeks computes the greeks of a european option according to Black Scholes formula
Black Scholes Greeks : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#The_Greeks
Possible values for Requests :  "all", "delta", "gamma", "vega", "theta", "rho"
Setting Request to "all" will compute all greeks
*/
func (s *europeanPricerServiceServerImpl) Compute(ctx context.Context, req *api_pb.ComputationRequest) (*api_pb.ComputationResponse, error) {
	//t := time.Unix(int64(req.Contract.Expiry), 0).Sub(
	//	time.Unix(int64(req.Pricingdate), 0)).Hours() / 24.0 / 365.250

	fmt.Print(req.Spot.Close)
	//bs(req.Spot.Close,
	//	req.Vol,
	//	req.Rate,
	//	req.Contract.Strike,
	//	t)

	return nil, nil
	//var (
	//	mult = putCallMap[strings.ToLower(in.Contract.Putcall)]
	//
	//	s = in.Marketdata.Spot.Index.Value
	//	v = in.Marketdata.Vol.Index.Value
	//	r = in.Marketdata.Rate.Index.Value
	//	k = in.Contract.Strike
	//	t = time.Unix(int64(in.Contract.Expiry), 0).Sub(
	//		time.Unix(int64(in.Pricingdate), 0)).Hours() / 24.0 / 365.250
	//)

	//if len(in.Greek) == 0 {
	//	return nil, errors.New("No greeks requested")
	//}
	//
	//sort.Strings(in.Greek)
	//if sort.SearchStrings(in.Greek, "all") < len(in.Greek) {
	//	in.Greek = allGreeks
	//}
	//
	//var (
	//	mult = putCallMap[strings.ToLower(in.Request.Contract.Putcall)]
	//
	//	s = in.Request.Marketdata.Spot.Index.Value
	//	v = in.Request.Marketdata.Vol.Index.Value
	//	r = in.Request.Marketdata.Rate.Index.Value
	//	k = in.Request.Contract.Strike
	//	t = time.Unix(int64(in.Request.Contract.Expiry), 0).Sub(
	//		time.Unix(int64(in.Request.Pricingdate), 0)).Hours() / 24.0 / 365.250
	//
	//	d1 = d1(s, k, t, v, r)
	//	d2 = d2(d1, v, t)
	//)
	//
	//greeks := make([]*pb.GreekResponse_Greek, len(in.Greek))
	//for index, greek := range in.Greek {
	//	greekResponse := &pb.GreekResponse_Greek{
	//		Label: greek,
	//	}
	//
	//	switch greek {
	//	case "delta":
	//		greekResponse.Value = delta(d1, mult)
	//	case "gamma":
	//		greekResponse.Value = gamma(s, t, v, d1)
	//	case "vega":
	//		greekResponse.Value = vega(s, t, d1)
	//	case "theta":
	//		greekResponse.Value = theta(s, k, t, v, r, d1, d2, mult)
	//	case "rho":
	//		greekResponse.Value = rho(k, t, r, d2, mult)
	//	default:
	//		greekResponse.Error = "Unknown greek " + greek
	//	}
	//
	//	greeks[index] = greekResponse
	//}
	//
	//return &pb.GreekResponse{
	//	Greeks: greeks,
	//}, nil

	return nil, nil
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
