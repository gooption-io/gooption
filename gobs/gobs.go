package gooption

import (
	"errors"
	"math"

	"github.com/ematvey/gostat"
)

var (
	phi  = stat.Normal_CDF(0, 1)
	dphi = stat.Normal_PDF(0, 1)
)

func BS(s, v, r, k, t, mult float64) float64 {
	d1 := D1(s, k, t, v, r)
	d2 := D2(d1, v, t)
	return mult * (s*phi(mult*d1) - k*phi(mult*d2)*math.Exp(-r*t))
}

type IVSolverResult struct {
	IV                float64
	NbSolverIteration int
	Error             error
}

func IVRootSolver(mktPrice, s, r, k, t, mult float64) *IVSolverResult {
	var (
		iv      = 0.1
		maxIter = 1000
	)

	for index := 0; index < maxIter; index++ {
		bsPrice := BS(s, iv, r, k, t, mult)
		iv = iv - (bsPrice-mktPrice)/Vega(s, t, D1(s, k, t, iv, r))
		if math.Abs(bsPrice-mktPrice) < 1E-10 { //decrease to 1E-25 to test convergence error
			return &IVSolverResult{iv, index, nil}
		}
	}
	return &IVSolverResult{iv, maxIter, errors.New("Did not converge to required interval")}
}

func BSGreek(label string, s, v, r, k, t, mult, d1, d2 float64) (float64, error) {
	var val float64
	switch label {
	case "delta":
		val = Delta(d1, mult)
	case "gamma":
		val = Gamma(s, t, v, d1)
	case "vega":
		val = Vega(s, t, d1)
	case "theta":
		val = Theta(s, k, t, v, r, d1, d2, mult)
	case "rho":
		val = Rho(k, t, r, d2, mult)
	default:
		return val, errors.New("Unknown greek " + label)
	}
	return val, nil
}

func D1(S, K, T, Sigma, R float64) float64 {
	return (1.0 / Sigma * math.Sqrt(T)) * (math.Log(S/K) + (R+Sigma*Sigma*0.5)*T)
}

func D2(d1, Sigma, T float64) float64 {
	return d1 - Sigma*math.Sqrt(T)
}

func Delta(d1, mult float64) float64 {
	return mult * phi(mult*d1)
}

func Gamma(s, t, sigma, d1 float64) float64 {
	return dphi(d1) / (s * sigma * math.Sqrt(t))
}

func Vega(s, t, d1 float64) float64 {
	return s * dphi(d1) * math.Sqrt(t)
}

func Theta(s, k, t, sigma, r, d1, d2, mult float64) float64 {
	return -0.5*(s*dphi(d1)*sigma/math.Sqrt(t)) - (mult * r * k * math.Exp(-r*t) * phi(mult*d2))
}

func Rho(k, t, r, d2, mult float64) float64 {
	return mult * k * t * math.Exp(-r*t) * phi(mult*d2)
}
