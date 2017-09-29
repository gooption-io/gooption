package main

import (
	"errors"
	"reflect"
	"sort"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/gooption/gobs/pb"
	"github.com/goyahoo"
)

var (
	generators = map[string]RequestGenerator{
		"greek":      greekRequestGenerator{},
		"impliedvol": impliedVolRequestGenerator{},
		"price":      priceRequestGenerator{},
	}

	pricingDate = float32(time.Now().Unix())
	option      = &pb.European{
		Ticker:  "AAPL DEC2017 PUT",
		Strike:  100.0,
		Expiry:  float32(time.Now().AddDate(0, 1, 0).Unix()),
		Putcall: pb.OptionType_CALL,
	}
	index = &pb.Index{
		Timestamp: pricingDate,
		Ticker:    "AAPL",
		Value:     100.0,
	}
	mkt = &pb.OptionMarket{
		Timestamp: pricingDate,
		Spot: &pb.Spot{
			Index: index,
		},
		Vol: &pb.FlatVol{
			Index: index,
		},
		Rate: &pb.RiskFreeRate{
			Index: index,
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
		Pricingdate: float32(time.Now().Unix()),
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
