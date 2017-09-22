package gooption

import (
	"math"

	"errors"

	"sort"

	"github.com/ematvey/gostat"
)

var (
	N         = stat.Normal_CDF(0, 1)
	n         = stat.Normal_PDF(0, 1)
	allGreeks = []string{"delta", "gamma", "vega", "theta", "rho"}
)

/*
FairValue computes the fair value of a european stock option according to Black Scholes formula
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Set PutCall to 1.0 for a Call option
Set PutCall to -1.0 for a Put option
Stock assumed to pay no dividends
*/
func FairValue(S, K, T, Sigma, R, PutCall float64) (float64, error) {
	if PutCall*PutCall != 1.0 {
		return 0, errors.New("putCall must be equal to 1(Call) or -1(Put)")
	}

	d1 := d1(S, K, T, Sigma, R)
	d2 := d2(d1, Sigma, T)
	price := PutCall * (S*N(PutCall*d1) - K*N(PutCall*d2)*math.Exp(-R*T))
	return price, nil
}

/*
Greeks computes the greeks of a european option according to Black Scholes formula
Black Scholes Greeks : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#The_Greeks
Set PutCall to 1.0 for a Call option
Set PutCall to -1.0 for a Put option
Possible values for Requests :  "all", "delta", "gamma", "vega", "theta", "rho"
Setting Request to "all" will compute all greeks
*/
func Greeks(S, K, T, Sigma, R, PutCall float64, Requests []string) ([]float64, []string, error) {
	if PutCall*PutCall != 1.0 {
		return nil, nil, errors.New("putCall must be equal to 1(Call) or -1(Put)")
	}
	if len(Requests) == 0 {
		return nil, nil, errors.New("No greek requested")
	}

	sort.Strings(Requests)
	if sort.SearchStrings(Requests, "all") < len(Requests) {
		Requests = allGreeks
	}

	greeks := make([]float64, len(Requests))
	greekErrors := make([]string, len(Requests))

	d1 := d1(S, K, T, Sigma, R)
	for index, actualIndex := 0, -1; index < len(Requests); index++ {
		if index > 0 && Requests[index] == Requests[index-1] {
			greeks[index] = greeks[actualIndex]
			continue
		}

		actualIndex++
		switch Requests[index] {
		case "delta":
			greeks[actualIndex] = delta(d1, PutCall)
		case "gamma":
			greeks[actualIndex] = gamma(S, T, Sigma, d1)
		case "vega":
			greeks[actualIndex] = vega(S, T, d1, PutCall)
		case "theta":
			greeks[actualIndex] = theta(S, K, T, Sigma, R, d1, d2(d1, Sigma, T), PutCall)
		case "rho":
			greeks[actualIndex] = rho(K, T, R, d2(d1, Sigma, T), PutCall)
		default:
			greekErrors[actualIndex] = "Unknown greek " + Requests[index]
		}
	}

	return greeks, greekErrors, nil
}

/*
ImpliedVol computes volatility matching the option quote passed as Quote using Newton Raphson solver
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
Set PutCall to 1.0 for a Call option
Set PutCall to -1.0 for a Put option
The second argument returned is the number of iteration used to converge
*/
func ImpliedVol(Quote, S, K, T, R, PutCall float64) (float64, int, error) {
	iv := 0.1
	maxIter := 1000
	for index := 0; index < maxIter; index++ {
		price, _ := FairValue(S, K, T, iv, R, PutCall)
		iv = iv - (price-Quote)/vega(S, T, d1(S, K, T, iv, R), PutCall)
		if math.Abs(price-Quote) < 1E-10 { //decrease to 1E-25 to test convergence error
			return iv, index, nil
		}
	}
	return iv, maxIter, errors.New("Did not converge to required interval")
}

func d1(S, K, T, Sigma, R float64) float64 {
	return (1.0 / Sigma * math.Sqrt(T)) * (math.Log(S/K) + (R+Sigma*Sigma*0.5)*T)
}

func d2(d1, Sigma, T float64) float64 {
	return d1 - Sigma*math.Sqrt(T)
}

func delta(d1, putCall float64) float64 {
	return putCall * N(putCall*d1)
}

func gamma(s, t, sigma, d1 float64) float64 {
	return n(d1) / (s * sigma * math.Sqrt(t))
}

func vega(s, t, d1, putCall float64) float64 {
	return s * n(d1) * math.Sqrt(t)
}

func theta(s, k, t, sigma, r, d1, d2, putCall float64) float64 {
	return -0.5*(s*n(d1)*sigma/math.Sqrt(t)) - (putCall * r * k * math.Exp(-r*t) * N(putCall*d2))
}

func rho(k, t, r, d2, putCall float64) float64 {
	return putCall * k * t * math.Exp(-r*t) * N(putCall*d2)
}
