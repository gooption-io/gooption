package main

import (
	"errors"
	"math"

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
func bs(s, v, r, k, t, mult float64) float64 {
	d1 := d1(s, k, t, v, r)
	d2 := d2(d1, v, t)

	return mult * (s*phi(mult*d1) - k*phi(mult*d2)*math.Exp(-r*t))
}

type ivSolverResult struct {
	IV                float64
	NbSolverIteration int
	Error             error
}

/*
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
*/
func ivRootSolver(mktPrice, s, r, k, t, mult float64) *ivSolverResult {
	var (
		iv      = 0.1
		maxIter = 1000
	)

	for index := 0; index < maxIter; index++ {
		bsPrice := bs(s, iv, r, k, t, mult)
		iv = iv - (bsPrice-mktPrice)/vega(s, t, d1(s, k, t, iv, r))
		if math.Abs(bsPrice-mktPrice) < 1E-10 { //decrease to 1E-25 to test convergence error
			return &ivSolverResult{iv, index, nil}
		}
	}
	return &ivSolverResult{iv, maxIter, errors.New("Did not converge to required interval")}
}

func bsGreek(label string, s, v, r, k, t, mult, d1, d2 float64) (float64, error) {
	var val float64
	switch label {
	case "delta":
		val = delta(d1, mult)
	case "gamma":
		val = gamma(s, t, v, d1)
	case "vega":
		val = vega(s, t, d1)
	case "theta":
		val = theta(s, k, t, v, r, d1, d2, mult)
	case "rho":
		val = rho(k, t, r, d2, mult)
	default:
		return val, errors.New("Unknown greek " + label)
	}
	return val, nil
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
