package main

import (
	"errors"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gooption/pb"
	"github.com/goyahoo"
)

var (
	dateFormat = "Jan 2, 2006"
)

type RequestGenerator interface {
	generate(ticker string) proto.Message
}

func NewRequest(request string) (proto.Message, error) {
	var generator RequestGenerator
	if strings.ToLower(request) == "price" {
		generator = priceRequestGenerator{}
	} else if strings.ToLower(request) == "greek" {
		generator = greekRequestGenerator{}
	} else if strings.ToLower(request) == "impliedvol" {
		generator = impliedVolRequestGenerator{}
	} else {
		return nil, errors.New("Unknown request type: " + request)
	}

	return generator.generate(ticker), nil
}

type impliedVolRequestGenerator struct{}

func (g impliedVolRequestGenerator) generate(ticker string) proto.Message {
	chain, _, _ := goyahoo.GetOptionChain(ticker)
	request := &pb.ImpliedVolRequest{
		Pricingdate: time.Now().Format(dateFormat),
		Quotes:      make([]*pb.ImpliedVolRequest_QuoteTermStructure, len(chain)),
		Marketdata: &pb.ImpliedVolRequest_MarketData{
			Spot:         100,
			Riskfreerate: 0.01,
		},
	}

	bindYahooQuotes := func(quotes []goyahoo.Quote) []*pb.ImpliedVolRequest_Quote {
		requestQuotes := make([]*pb.ImpliedVolRequest_Quote, len(quotes))
		for i, quote := range quotes {
			requestQuotes[i] = &pb.ImpliedVolRequest_Quote{
				Strike:       quote.Strike,
				Ask:          quote.Ask,
				Bid:          quote.Bid,
				Openinterest: int32(quote.OpenInterest),
			}
		}
		return requestQuotes
	}

	for i, yahooQuote := range chain {
		request.Quotes[i] = &pb.ImpliedVolRequest_QuoteTermStructure{
			Expiry: float64(yahooQuote.ExpirationDates[i]),
			Puts:   bindYahooQuotes(yahooQuote.Options[0].Puts),
			Calls:  bindYahooQuotes(yahooQuote.Options[0].Calls),
		}
	}

	return request
}

type priceRequestGenerator struct{}

func (g priceRequestGenerator) generate(ticker string) proto.Message {
	return &pb.PriceRequest{
		Pricingdate: "Jan 2, 2006",
		Contracts: []*pb.Contract{
			{
				Strike:  100.0,
				Expiry:  "Jan 2, 2007",
				Putcall: pb.Contract_CALL,
			},
		},
		Marketdata: []*pb.PriceRequest_MarketData{
			{
				Spot:         100,
				Volatility:   0.10,
				Riskfreerate: 0.01,
			},
		},
	}
}

type greekRequestGenerator struct{}

func (g greekRequestGenerator) generate(ticker string) proto.Message {
	return &pb.GreekRequest{
		Request: &pb.PriceRequest{
			Pricingdate: "Jan 2, 2006",
			Contracts: []*pb.Contract{
				{
					Strike:  100.0,
					Expiry:  "Jan 2, 2007",
					Putcall: pb.Contract_CALL,
				},
			},
			Marketdata: []*pb.PriceRequest_MarketData{
				{
					Spot:         100,
					Volatility:   0.10,
					Riskfreerate: 0.01,
				},
			},
		},
		Greek: []string{"all"},
	}
}
