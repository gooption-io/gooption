// -- //go:generate sh -c "protoc --proto_path=pb --proto_path=$GOPATH/src/github.com/gooption/pb --proto_path=$GOPATH/src --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogofast_out=plugins=grpc:pb $GOPATH/src/github.com/gooption/pb/*.proto"
// go:generate "protoc.exe --proto_path=pb --proto_path=%GOPATH%\\src\\github.com\\lehajam\\gooption\\pb --proto_path=%GOPATH%\\src --proto_path=%GOPATH%\\src\\github.com\\grpc-ecosystem\\grpc-gateway\\third_party\\googleapis --gogofast_out=plugins=grpc:pb %GOPATH%\\src\\github.com\\gooption\\lehajam\\pb\\*.proto"
package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/lehajam/gooption/gooption-cli/pb"
	"github.com/lehajam/goyahoo"
)

var (
	generators = map[string]RequestGenerator{
		"greek":      greekRequestGenerator{},
		"impliedvol": impliedVolRequestGenerator{},
		"price":      priceRequestGenerator{},
	}

	pricingDate = float64(time.Now().Unix())
	option      = &pb.European{
		Timestamp: pricingDate,
		Ticker:    "AAPL DEC2017 PUT",
		Undticker: "AAPL",
		Strike:    159.76,
		Expiry:    float64(time.Now().AddDate(0, 1, 0).Unix()),
		Putcall:   "put",
	}
	mkt = &pb.OptionMarket{
		Timestamp: pricingDate,
		Spot: &pb.Spot{
			Index: &pb.Index{
				Timestamp: pricingDate,
				Ticker:    "AAPL",
				Value:     159.76,
			},
		},
		Vol: &pb.FlatVol{
			Index: &pb.Index{
				Timestamp: pricingDate,
				Ticker:    "AAPL",
				Value:     0.10,
			},
		},
		Rate: &pb.RiskFreeRate{
			Index: &pb.Index{
				Timestamp: pricingDate,
				Ticker:    "USD.FEDFUND",
				Value:     0.01,
			},
		},
	}
)

type RequestGenerator interface {
	generate(ticker string) (proto.Message, error)
}

func NewRequest(request string) (proto.Message, error) {
	keys := reflect.ValueOf(generators).MapKeys()
	if sort.Search(len(keys), func(i int) bool { return keys[i].String() == request }) > len(keys) {
		return nil, errors.New("Unknown request type: " + request)
	}

	return generators[request].generate(ticker)
}

type priceRequestGenerator struct{}

func (g priceRequestGenerator) generate(ticker string) (proto.Message, error) {
	return &pb.PriceRequest{
		Pricingdate: pricingDate,
		Contract:    option,
		Marketdata:  mkt,
	}, nil
}

type greekRequestGenerator struct{}

func (g greekRequestGenerator) generate(ticker string) (proto.Message, error) {
	return &pb.GreekRequest{
		Request: &pb.PriceRequest{
			Pricingdate: pricingDate,
			Contract:    option,
			Marketdata:  mkt,
		},
		Greek: []string{"all"},
	}, nil
}

type impliedVolRequestGenerator struct{}

func (g impliedVolRequestGenerator) bind(quotes []goyahoo.Quote, putcall string) []*pb.OptionQuote {
	requestQuotes := make([]*pb.OptionQuote, len(quotes))
	for i, quote := range quotes {
		requestQuotes[i] = &pb.OptionQuote{
			Timestamp:    pricingDate,
			Strike:       quote.Strike,
			Ask:          quote.Ask,
			Bid:          quote.Bid,
			Openinterest: quote.OpenInterest,
			Putcall:      putcall,
		}
	}
	return requestQuotes
}

func (g impliedVolRequestGenerator) generate(ticker string) (proto.Message, error) {
	chain, _, err := goyahoo.GetOptionChain(ticker)
	if err != nil {
		return nil, err
	}

	request := &pb.ImpliedVolRequest{
		Pricingdate: pricingDate,
		Marketdata:  mkt,
		Quotes:      make([]*pb.OptionQuoteSlice, len(chain)),
	}

	for i, yahooQuote := range chain {
		request.Quotes[i] = &pb.OptionQuoteSlice{
			Timestamp: pricingDate,
			Expiry:    float64(yahooQuote.ExpirationDates[i]),
			Puts:      g.bind(yahooQuote.Options[0].Puts, "put"),
			Calls:     g.bind(yahooQuote.Options[0].Calls, "call"),
		}
	}

	fmt.Println(request)
	return request, nil
}
