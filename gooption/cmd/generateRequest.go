// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	gooption "github.com/gooption-io/gooption/proto"
	"github.com/lehajam/goyahoo"
	"github.com/spf13/cobra"
)

var (
	requestType                                  string
	strike, spot, vol, rate                      float64
	optionTicker, undTicker, rateTicker, putcall string
	pricingDate                                  = float64(time.Now().Unix())

	option = &gooption.European{
		Timestamp: pricingDate,
		Ticker:    optionTicker,
		Undticker: undTicker,
		Strike:    strike,
		Expiry:    float64(time.Now().AddDate(0, 1, 0).Unix()),
		Putcall:   putcall,
	}

	mkt = &gooption.OptionMarket{
		Timestamp: pricingDate,
		Spot: &gooption.Spot{
			Index: &gooption.Index{
				Timestamp: pricingDate,
				Ticker:    undTicker,
				Value:     spot,
			},
		},
		Vol: &gooption.FlatVol{
			Index: &gooption.Index{
				Timestamp: pricingDate,
				Ticker:    undTicker,
				Value:     vol,
			},
		},
		Rate: &gooption.RiskFreeRate{
			Index: &gooption.Index{
				Timestamp: pricingDate,
				Ticker:    rateTicker,
				Value:     rate,
			},
		},
	}

	priceRequest = &gooption.PriceRequest{
		Pricingdate: pricingDate,
		Contract:    option,
		Marketdata:  mkt,
	}

	greekRequest = &gooption.GreekRequest{
		Greek: []string{"all"},
		Request: &gooption.PriceRequest{
			Pricingdate: pricingDate,
			Contract:    option,
			Marketdata:  mkt,
		},
	}
)

func init() {
	generateCmd.AddCommand(generateRequestCmd)

	generateCmd.Flags().StringVar(&requestType, "type", "", "request type eg. price greek or impliedvol. All if not specified")
	generateCmd.Flags().StringVar(&undTicker, "und", "AAPL", "underlying ticker")
	generateCmd.Flags().StringVar(&putcall, "putcall", "put", "option type (call or put")
	generateCmd.Flags().StringVar(&optionTicker, "ticker", "AAPL DEC2017 PUT", "option ticker")
	generateCmd.Flags().StringVar(&rateTicker, "rateticker", "AAPL DEC2017 PUT", "option ticker")
	generateCmd.Flags().Float64Var(&vol, "vol", 0.10, "vol level")
	generateCmd.Flags().Float64Var(&rate, "rate", 0.01, "rate level")
	generateCmd.Flags().Float64Var(&spot, "spot", 159.76, "spot level")
}

// generateCmd represents the generate command
var generateRequestCmd = &cobra.Command{
	Use:   "request",
	Short: "Generate mock request",
	Long: `Generate json files containing mock requests in current folder.
Generate one request per file, see help for available flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch requestType {
		case "price":
			generateRequest("pricerequest.json", priceRequest)
		case "greek":
			generateRequest("greekrequest.json", greekRequest)
		case "impliedvol":
			generateRequest("impliedvol.json", newImpliedVolRequest())
		default:
			generateRequest("pricerequest.json", priceRequest)
			generateRequest("greekrequest.json", greekRequest)
			generateRequest("impliedvol.json", newImpliedVolRequest())
		}
	},
}

func generateRequest(name string, message proto.Message) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	marshaler := jsonpb.Marshaler{EnumsAsInts: true, EmitDefaults: true, Indent: "    "}
	err = marshaler.Marshal(f, message)
	if err != nil {
		panic(err)
	}
}

func newImpliedVolRequest() *gooption.ImpliedVolRequest {
	chain, _, err := goyahoo.GetOptionChain(undTicker)
	if err != nil {
		panic(err)
	}

	request := &gooption.ImpliedVolRequest{
		Pricingdate: pricingDate,
		Marketdata:  mkt,
		Quotes:      make([]*gooption.OptionQuoteSlice, len(chain)),
	}

	for i, yahooQuote := range chain {
		request.Quotes[i] = &gooption.OptionQuoteSlice{
			Timestamp: pricingDate,
			Expiry:    float64(yahooQuote.ExpirationDates[i]),
			Puts:      newOptionQuote(yahooQuote.Options[0].Puts, "put"),
			Calls:     newOptionQuote(yahooQuote.Options[0].Calls, "call"),
		}
	}

	return request
}

func newOptionQuote(quotes []goyahoo.Quote, putcall string) []*gooption.OptionQuote {
	requestQuotes := make([]*gooption.OptionQuote, len(quotes))
	for i, quote := range quotes {
		requestQuotes[i] = &gooption.OptionQuote{
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
