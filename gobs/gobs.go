package main

import (
	"math"
	"strings"
	"time"

	"github.com/gooption-io/gooption/gobs/pb"
	"gonum.org/v1/gonum/stat/distuv"
)

var (
	phi  = distuv.Normal{Mu: 0, Sigma: 1}.CDF
	dphi = distuv.Normal{Mu: 0, Sigma: 1}.Prob

	allGreeks  = []string{"delta", "gamma", "vega", "theta", "rho"}
	putCallMap = map[string]float64{
		"call": 1.0,
		"put":  -1.0,
	}
)

/*
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
*/
func bs(pricingDate float64, contract *pb.European, mkt *pb.OptionMarket) float64 {
	var (
		mult = putCallMap[strings.ToLower(contract.Putcall)]

		s = mkt.Spot.Index.Value
		v = mkt.Vol.Index.Value
		r = mkt.Rate.Index.Value
		k = contract.Strike
		t = time.Unix(int64(contract.Expiry), 0).Sub(
			time.Unix(int64(pricingDate), 0)).Hours() / 24.0 / 365.250
	)

	d1 := d1(s, k, t, v, r)
	d2 := d2(d1, v, t)

	return mult * (s*phi(mult*d1) - k*phi(mult*d2)*math.Exp(-r*t))
}

func bsGreek(in *pb.GreekRequest) []*pb.GreekResponse_Greek {
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

	return greeks
}

func bsImpliedVol(index int, in *pb.ImpliedVolRequest, out chan<- pb.ImpliedVolSlice) {
	var (
		t = func(q *pb.OptionQuoteSlice) float64 {
			return time.Unix(int64(q.Expiry), 0).Sub(time.Unix(int64(in.Pricingdate), 0)).Hours() / 24.0 / 365.250
		}
	)

	slice := in.Quotes[index]
	spot := in.Marketdata.Spot.Index.Value
	calibratedSlice := pb.ImpliedVolSlice{
		Timestamp: in.Marketdata.Timestamp,
		Expiry:    slice.Expiry,
		Iserror:   false,
		Quotes:    make([]*pb.ImpliedVolQuote, len(slice.Puts)+len(slice.Calls)),
	}

	calibIndex := 0
	for _, put := range slice.Puts {
		if put.Strike/spot > putLBound && put.Strike/spot <= 1.0 {
			calibratedSlice.Quotes[calibIndex] = ivRootSolver(in.Pricingdate, t(slice), put, in.Marketdata)
			if calibratedSlice.Quotes[calibIndex].Error != "" {
				calibratedSlice.Iserror = true
			}
			calibIndex++
		}
	}

	for _, call := range slice.Calls {
		if call.Strike/spot > 1.0 {
			calibratedSlice.Quotes[calibIndex] = ivRootSolver(in.Pricingdate, t(slice), call, in.Marketdata)
			if calibratedSlice.Quotes[calibIndex].Error != "" {
				calibratedSlice.Iserror = true
			}
			calibIndex++
		}
	}

	calibratedSlice.Quotes = calibratedSlice.Quotes[0:calibIndex]
	out <- calibratedSlice
}

/*
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
*/
func ivRootSolver(pricingDate, expiry float64, quote *pb.OptionQuote, mkt *pb.OptionMarket) *pb.ImpliedVolQuote {
	const (
		iv      = 0.1
		maxIter = 1000
	)

	k := quote.Strike
	s := mkt.Spot.Index.Value
	r := mkt.Rate.Index.Value
	mktPrice := (quote.Ask + quote.Bid) / 2.0
	contract := &pb.European{Strike: k, Putcall: quote.Putcall, Expiry: expiry}

	for index := 0; index < maxIter; index++ {
		bsPrice := bs(pricingDate, contract, mkt)
		iv := iv - (bsPrice-mktPrice)/vega(s, expiry, d1(s, k, expiry, iv, r))
		if math.Abs(bsPrice-mktPrice) < 1E-10 { //decrease to 1E-25 to test convergence error
			return &pb.ImpliedVolQuote{
				Input:       quote,
				Vol:         iv,
				Nbiteration: int64(index),
			}
		}
	}

	return &pb.ImpliedVolQuote{
		Input:       quote,
		Error:       "Did not converge to required interval",
		Nbiteration: maxIter,
	}
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
