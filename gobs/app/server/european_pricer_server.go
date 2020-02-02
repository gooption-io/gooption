package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"gonum.org/v1/gonum/stat/distuv"

	api_pb "github.com/gooption-io/gooption/v1/gobs/api"

	"github.com/izumin5210/grapi/pkg/grapiserver"
)

const (
	ivSeed    = 0.1 // solver starting point
	maxIter   = 1000
	putLBound = 0.20
)

var (
	phi             = distuv.Normal{Mu: 0, Sigma: 1}.CDF
	dphi            = distuv.Normal{Mu: 0, Sigma: 1}.Prob
	mapToMultiplier = map[string]float64{"call": 1.0, "put": -1.0}
	allGreeks       = []string{"delta", "gamma", "vega", "theta", "rho"}
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
func (srv *europeanPricerServiceServerImpl) Compute(ctx context.Context, req *api_pb.ComputationRequest) (*api_pb.ComputationResponse, error) {
	var (
		s = req.Spot.Close
		v = req.Vol
		r = req.Rate
		k = req.Contract.Strike
		t = time.Unix(int64(req.Contract.Expiry), 0).Sub(
			time.Unix(int64(req.Pricingdate), 0)).Hours() / 24.0 / 365.250

		mult = mapToMultiplier[strings.ToLower(req.Contract.Putcall)]
	)

	jreq, _ := json.MarshalIndent(req, "", "\t")
	fmt.Printf("%s \n", jreq)

	sort.Strings(req.Greek)
	if sort.SearchStrings(req.Greek, "all") < len(req.Greek) {
		req.Greek = allGreeks
	}

	d1 := d1(s, k, t, v, r)
	d2 := d2(d1, v, t)

	greeks := make([]*api_pb.Greek, len(req.Greek))
	for index, name := range req.Greek {
		greek := &api_pb.Greek{
			Label: name,
		}

		switch name {
		case "delta":
			greek.Value = delta(d1, mult)
		case "gamma":
			greek.Value = gamma(s, t, v, d1)
		case "vega":
			greek.Value = vega(s, t, d1)
		case "theta":
			greek.Value = theta(s, k, t, v, r, d1, d2, mult)
		case "rho":
			greek.Value = rho(k, t, r, d2, mult)
		default:
			greek.Error = "Unknown greek " + name
		}

		greeks[index] = greek
	}

	return &api_pb.ComputationResponse{
		Price:  bs(s, v, r, k, t, mult),
		Greeks: greeks,
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
