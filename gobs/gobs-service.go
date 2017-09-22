//go:generate gooption-gen gobs Price Greek ImpliedVol
package gooption

import (
	"time"

	"github.com/gooption/pb"
)

var (
	dateFormat = "Jan 2, 2006"
	putCallMap = map[pb.Contract_OptionType]float64{
		pb.Contract_CALL: 1.0,
		pb.Contract_PUT:  -1.0,
	}
)

type Service struct {
}

func (b Service) Price(request *pb.PriceRequest) (*pb.PriceResponse, error) {
	response := &pb.PriceResponse{
		Price: make([]float64, len(request.Contracts)),
		Error: make([]string, len(request.Contracts)),
	}

	pricingDate, err := time.Parse(dateFormat, request.Pricingdate)
	if err != nil {
		return response, err
	}

	for index := 0; index < len(request.Contracts); index++ {
		expiry, err := time.Parse(dateFormat, request.Contracts[index].Expiry)
		if err != nil {
			response.Error[index] = err.Error()
			continue
		}

		fv, err := FairValue(
			request.Marketdata[index].Spot,
			request.Contracts[index].Strike,
			expiry.Sub(pricingDate).Hours()/24.0/365.25,
			request.Marketdata[index].Volatility,
			request.Marketdata[index].Riskfreerate,
			putCallMap[request.Contracts[index].Putcall])

		response.Price[index] = fv
		if err != nil {
			response.Error[index] = err.Error()
		}
	}

	return response, nil
}

func (b Service) Greek(request *pb.GreekRequest) (*pb.GreekResponse, error) {
	priceRequest := request.Request
	response := &pb.GreekResponse{
		Greeks: make([]*pb.GreekResponse_Greek, len(priceRequest.Contracts)),
	}

	pricingDate, err := time.Parse(dateFormat, priceRequest.Pricingdate)
	if err != nil {
		return nil, err
	}
	for index := 0; index < len(request.Request.Contracts); index++ {
		expiry, err := time.Parse(dateFormat, priceRequest.Contracts[index].Expiry)
		if err != nil {
			response.Errors[index] = err.Error()
			continue
		}

		greeks, greekErrors, err := Greeks(
			priceRequest.Marketdata[index].Spot,
			priceRequest.Contracts[index].Strike,
			expiry.Sub(pricingDate).Hours()/24.0/365.25,
			priceRequest.Marketdata[index].Volatility,
			priceRequest.Marketdata[index].Riskfreerate,
			putCallMap[priceRequest.Contracts[index].Putcall],
			request.Greek)

		if err != nil {
			response.Errors[index] = err.Error()
		} else {
			response.Greeks[index] = &pb.GreekResponse_Greek{
				Labels: request.Greek,
				Values: greeks,
				Errors: greekErrors,
			}
		}
	}

	return response, nil
}

func (b Service) ImpliedVol(request *pb.ImpliedVolRequest) (*pb.ImpliedVolResponse, error) {
	pricingDate, err := time.Parse(dateFormat, request.Pricingdate)
	if err != nil {
		return nil, err
	}

	volSurface, err := buildVolSurface(
		pricingDate,
		request.Marketdata.Spot,
		request.Marketdata.Riskfreerate,
		request.Quotes)
	if err != nil {
		return nil, err
	}

	response := &pb.ImpliedVolResponse{
		Volsurface: make([]*pb.ImpliedVolResponse_ImpliedVolTermStructure, len(volSurface)),
	}

	for i, slice := range volSurface {
		response.Volsurface[i] = &pb.ImpliedVolResponse_ImpliedVolTermStructure{
			Iserror: slice.IsError,
			Expiry:  slice.Expiry,
			Strikes: slice.Strikes,
			Vols:    slice.Values,
			Errors:  slice.Errors,
		}
	}

	return response, nil
}
