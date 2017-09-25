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

type BSService struct {
	Service
}

func NewService() Service {
	return BSService{}
}

/*
Price computes the fair value of a european stock option according to Black Scholes formula
Black Scholes Formula : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#Black.E2.80.93Scholes_formula
Stock assumed to pay no dividends
*/
func (service BSService) Price(request *pb.PriceRequest) *pb.PriceResponse {
	var (
		mult = putCallMap[request.Contract.Putcall]

		s = request.Marketdata.Spot.Value
		v = request.Marketdata.Vol.Value
		r = request.Marketdata.Rate.Value
		k = request.Contract.Strike
		t = time.Unix(int64(request.Contract.Expiry), 0).Sub(
			time.Unix(int64(request.Pricingdate), 0)).Hours() / 24.0 / 365.250
		bs = BS(s, v, r, k, t, mult)
	)

	return &pb.PriceResponse{
		Price: bs,
	}
}

/*
Greeks computes the greeks of a european option according to Black Scholes formula
Black Scholes Greeks : https://en.wikipedia.org/wiki/Black%E2%80%93Scholes_model#The_Greeks
Set PutCall to 1.0 for a Call option
Set PutCall to -1.0 for a Put option
Possible values for Requests :  "all", "delta", "gamma", "vega", "theta", "rho"
Setting Request to "all" will compute all greeks
*/
func (service BSService) Greek(request *pb.GreekRequest) *pb.GreekResponse {
	if len(request.Greek) == 0 {
		return &pb.GreekResponse{Error: "No greeks requested"}
	}

	sort.Strings(request.Greek)
	if sort.SearchStrings(request.Greek, "all") < len(request.Greek) {
		request.Greek = allGreeks
	}

	return newGreekResponse(request)
}

func newGreekResponse(request *pb.GreekRequest) *pb.GreekResponse {
	var (
		mult = putCallMap[request.Request.Contract.Putcall]

		s = request.Request.Marketdata.Spot.Value
		v = request.Request.Marketdata.Vol.Value
		r = request.Request.Marketdata.Rate.Value
		k = request.Request.Contract.Strike
		t = time.Unix(int64(request.Request.Contract.Expiry), 0).Sub(
			time.Unix(int64(request.Request.Pricingdate), 0)).Hours() / 24.0 / 365.250
		d1 = D1(s, k, t, v, r)
		d2 = D2(d1, v, t)
	)

	newGreek := func(label string) *pb.GreekResponse_Greek {
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
			return &pb.GreekResponse_Greek{
				Label: label,
				Error: "Unknown greek " + label,
			}
		}

		return &pb.GreekResponse_Greek{
			Label: label,
			Value: val,
		}
	}

	response := &pb.GreekResponse{
		Greeks: make([]*pb.GreekResponse_Greek, len(request.Greek)),
	}

	for index := 0; index < len(request.Greek); index++ {
		response.Greeks[index] = newGreek(request.Greek[index])
	}

	return response
}

/*
ImpliedVol computes volatility matching the option quote passed as Quote using Newton Raphson solver
Newton Raphson solver : https://en.wikipedia.org/wiki/Newton%27s_method
Set PutCall to 1.0 for a Call option
Set PutCall to -1.0 for a Put option
The second argument returned is the number of iteration used to converge
*/
func (service BSService) ImpliedVol(request *pb.ImpliedVolRequest) *pb.ImpliedVolResponse {
	surf := &pb.ImpliedVolSurface{
		Timestamp:  request.Marketdata.Timestamp,
		Volsurface: make([]*pb.ImpliedVolSlice, len(request.Quotes)),
	}

	for index := 0; index < len(request.Quotes); index++ {
		surf.Volsurface[index] = NewImpliedVolSlice(int64(request.Pricingdate), request.Quotes[index], request.Marketdata)
	}

	return &pb.ImpliedVolResponse{
		Volsurface: surf,
	}
}
