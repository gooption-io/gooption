//go:generate gooption-cli -p gobs -r Price -r Greek -r ImpliedVol
package gooption

import (
	"sort"
	"time"

	"github.com/gooption/pb"
)

var (
	allGreeks  = []string{"delta", "gamma", "vega", "theta", "rho"}
	putCallMap = map[pb.OptionType]float64{
		pb.OptionType_CALL: 1.0,
		pb.OptionType_PUT:  -1.0,
	}
)

type Service interface {
	Price(request *pb.PriceRequest) *pb.PriceResponse
	Greek(request *pb.GreekRequest) *pb.GreekResponse
	ImpliedVol(request *pb.ImpliedVolRequest) *pb.ImpliedVolResponse
}

func NewService() Service {
	return gobsService{}
}

type gobsService struct {
	Service
}

/*
Price computes the fair value of a european stock option according to Black Scholes formula
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
*/
func (service gobsService) Price(request *pb.PriceRequest) *pb.PriceResponse {
	var (
		mult = putCallMap[request.Contract.Putcall]

		s = request.Marketdata.Spot.Value
		v = request.Marketdata.Vol.Value
		r = request.Marketdata.Rate.Value
		k = request.Contract.Strike
		t = time.Unix(int64(request.Contract.Expiry), 0).Sub(
			time.Unix(int64(request.Pricingdate), 0)).Hours() / 24.0 / 365.250
		bs = bs(s, v, r, k, t, mult)
	)

	return &pb.PriceResponse{
		Price: bs,
	}
}

/*
Greeks computes the greeks of a european option according to Black Scholes formula
Black Scholes Greeks : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#The_Greeks
Possible values for Requests :  "all", "delta", "gamma", "vega", "theta", "rho"
Setting Request to "all" will compute all greeks
*/
func (service gobsService) Greek(request *pb.GreekRequest) *pb.GreekResponse {
	var (
		mult = putCallMap[request.Request.Contract.Putcall]

		s = request.Request.Marketdata.Spot.Value
		v = request.Request.Marketdata.Vol.Value
		r = request.Request.Marketdata.Rate.Value
		k = request.Request.Contract.Strike
		t = time.Unix(int64(request.Request.Contract.Expiry), 0).Sub(
			time.Unix(int64(request.Request.Pricingdate), 0)).Hours() / 24.0 / 365.250
		d1 = d1(s, k, t, v, r)
		d2 = d2(d1, v, t)
	)

	if len(request.Greek) == 0 {
		return &pb.GreekResponse{Error: "No greeks requested"}
	}

	sort.Strings(request.Greek)
	if sort.SearchStrings(request.Greek, "all") < len(request.Greek) {
		request.Greek = allGreeks
	}

	response := &pb.GreekResponse{
		Greeks: make([]*pb.GreekResponse_Greek, len(request.Greek)),
	}

	for index := 0; index < len(request.Greek); index++ {
		response.Greeks[index] = &pb.GreekResponse_Greek{
			Label: request.Greek[index],
		}

		greek, err := bsGreek(response.Greeks[index].Label, s, v, r, k, t, mult, d1, d2)
		if err != nil {
			response.Greeks[index].Error = err.Error()
		}
		response.Greeks[index].Value = greek
	}

	return response
}

/*
ImpliedVol computes volatility matching the option quote passed as Quote using Newton Raphson solver
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
The second argument returned is the number of iteration used to converge
*/
func (service gobsService) ImpliedVol(request *pb.ImpliedVolRequest) *pb.ImpliedVolResponse {
	var (
		mult = func(q *pb.OptionQuote) float64 { return putCallMap[q.Putcall] }

		s = request.Marketdata.Spot.Value
		r = request.Marketdata.Rate.Value
		k = func(q *pb.OptionQuote) float64 { return q.Strike }
		p = func(q *pb.OptionQuote) float64 { return (q.Ask + q.Bid) / 2.0 }
		t = func(q *pb.OptionQuoteSlice) float64 {
			return time.Unix(int64(q.Expiry), 0).Sub(time.Unix(int64(request.Pricingdate), 0)).Hours() / 24.0 / 365.250
		}
	)

	surf := &pb.ImpliedVolSurface{
		Timestamp:  request.Marketdata.Timestamp,
		Volsurface: make([]*pb.ImpliedVolSlice, len(request.Quotes)),
	}

	for index := 0; index < len(request.Quotes); index++ {
		newOptionQuoteSliceIterator(request.Quotes[index], request.Marketdata).foreach(
			func(quote *pb.OptionQuote) *ivSolverResult {
				return ivRootSolver(p(quote), s, r, k(quote), t(request.Quotes[index]), mult(quote))
			}).then(
			func(calibratedSlice *pb.ImpliedVolSlice) {
				surf.Volsurface[index] = calibratedSlice
			})
	}

	return &pb.ImpliedVolResponse{
		Volsurface: surf,
	}
}
