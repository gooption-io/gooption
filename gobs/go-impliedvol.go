package gooption

import (
	"time"

	"github.com/gooption/pb"
)

func NewImpliedVolSlice(pricingDate int64, quotes *pb.OptionQuoteSlice, market *pb.OptionMarket) *pb.ImpliedVolSlice {
	var (
		mult = func(q *pb.OptionQuote) float64 { return putCallMap[q.Putcall] }

		s = market.Spot.Value
		r = market.Rate.Value
		k = func(q *pb.OptionQuote) float64 { return q.Strike }
		p = func(q *pb.OptionQuote) float64 { return (q.Ask + q.Bid) / 2.0 }
		t = time.Unix(int64(quotes.Expiry), 0).Sub(time.Unix(int64(pricingDate), 0)).Hours() / 24.0 / 365.250
		n = len(quotes.Puts) + len(quotes.Calls)
	)

	count, isError, errors := 0, false, make([]string, n)
	calibQuotes, vols, strikes := make([]*pb.OptionQuote, n), make([]float64, n), make([]float64, n)
	optionQuoteSliceIterator{quotes, market}.foreach(
		func(index int, quote *pb.OptionQuote) {
			if quote != nil {
				count++
				calibQuotes[index] = quote
				if iv, _, err := IVRootSolver(p(quote), s, r, k(quote), t, mult(quote)); err == nil {
					strikes[index], vols[index] = k(quote), iv
				} else {
					isError, errors[index] = true, err.Error()
				}
			}
		})

	return &pb.ImpliedVolSlice{
		Timestamp: market.Timestamp,
		Expiry:    quotes.Expiry,
		Iserror:   isError,
		Vols:      vols[0:count],
		Strikes:   strikes[0:count],
		Errors:    errors[0:count],
		Quotes:    calibQuotes[0:count],
	}
}

type optionQuoteSliceIterator struct {
	Slice  *pb.OptionQuoteSlice
	Market *pb.OptionMarket
}

func (it optionQuoteSliceIterator) foreach(f func(idx int, quote *pb.OptionQuote)) {
	count := 0
	for index := 0; index < len(it.Slice.Puts); index++ {
		if index < len(it.Slice.Puts) && it.Slice.Puts[index].Strike <= it.Market.Spot.Value {
			if it.Slice.Puts[index].Strike/it.Market.Spot.Value > 0.20 {
				f(count, it.Slice.Puts[index])
				count++
			}
		}
	}

	for index := 0; index < len(it.Slice.Calls); index++ {
		if index < len(it.Slice.Calls) && it.Slice.Calls[index].Strike > it.Market.Spot.Value {
			f(count, it.Slice.Calls[index])
			count++
		}
	}
}
