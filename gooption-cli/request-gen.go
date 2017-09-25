package main

import (
	"errors"
	"reflect"
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/gooption/pb"
	"github.com/goyahoo"
)

var (
	generators = map[string]RequestGenerator{
		"greek":      greekRequestGenerator{},
		"impliedvol": impliedVolRequestGenerator{},
		"price":      priceRequestGenerator{},
	}

	pricingDate float32 = 1136192400
	option              = &pb.Contract{
		Strike:  100.0,
		Expiry:  1167728400,
		Putcall: pb.OptionType_CALL,
	}
	mkt = &pb.OptionMarket{
		Timestamp: pricingDate,
		Spot: &pb.Spot{
			Timestamp: pricingDate,
			Ticker:    "AAPL",
			Value:     100.0,
		},
		Vol: &pb.FlatVol{
			Timestamp: pricingDate,
			Ticker:    "AAPL",
			Value:     0.10,
		},
		Rate: &pb.RiskFreeRate{
			Timestamp: pricingDate,
			Ticker:    "USD.FEDFUND",
			Value:     0.01,
		},
	}
)

type RequestGenerator interface {
	generate(ticker string) proto.Message
}

func NewRequest(request string) (proto.Message, error) {
	keys := reflect.ValueOf(generators).MapKeys()
	if sort.Search(len(keys), func(i int) bool { return keys[i].String() == request }) > len(keys) {
		return nil, errors.New("Unknown request type: " + request)
	}

	return generators[request].generate(ticker), nil
}

type priceRequestGenerator struct{}
type greekRequestGenerator struct{}
type impliedVolRequestGenerator struct{}

func (g priceRequestGenerator) generate(ticker string) proto.Message {
	return &pb.PriceRequest{
		Pricingdate: pricingDate,
		Contract:    option,
		Marketdata:  mkt,
	}
}

func (g greekRequestGenerator) generate(ticker string) proto.Message {

	return &pb.GreekRequest{
		Request: &pb.PriceRequest{
			Pricingdate: pricingDate,
			Contract:    option,
			Marketdata:  mkt,
		},
		Greek: []string{"all"},
	}
}

func (g impliedVolRequestGenerator) bind(quotes []goyahoo.Quote, putcall pb.OptionType) []*pb.OptionQuote {
	requestQuotes := make([]*pb.OptionQuote, len(quotes))
	for i, quote := range quotes {
		requestQuotes[i] = &pb.OptionQuote{
			Strike:       quote.Strike,
			Ask:          quote.Ask,
			Bid:          quote.Bid,
			Openinterest: float32(quote.OpenInterest),
			Putcall:      putcall,
		}
	}
	return requestQuotes
}

func (g impliedVolRequestGenerator) generate(ticker string) proto.Message {
	chain, _, _ := goyahoo.GetOptionChain(ticker)
	request := &pb.ImpliedVolRequest{
		Pricingdate: pricingDate,
		Marketdata:  mkt,
		Quotes:      make([]*pb.OptionQuoteSlice, len(chain)),
	}

	for i, yahooQuote := range chain {
		request.Quotes[i] = &pb.OptionQuoteSlice{
			Timestamp: request.Pricingdate,
			Expiry:    float32(yahooQuote.ExpirationDates[i]),
			Puts:      g.bind(yahooQuote.Options[0].Puts, pb.OptionType_PUT),
			Calls:     g.bind(yahooQuote.Options[0].Calls, pb.OptionType_CALL),
		}
	}

	return request
}
