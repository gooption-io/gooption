package main

import (
	"github.com/gooption/gobs/pb"
)

var (
	putLBound = 0.20
)

type optionQuoteSliceIterator struct {
	NbCalibratedSlice int
	Slice             *pb.OptionQuoteSlice
	Market            *pb.OptionMarket
	CalibratedSlice   *pb.ImpliedVolSlice
}

func newOptionQuoteSliceIterator(quotes *pb.OptionQuoteSlice, market *pb.OptionMarket) *optionQuoteSliceIterator {
	var (
		n = len(quotes.Puts) + len(quotes.Calls)
	)

	return &optionQuoteSliceIterator{
		Slice:  quotes,
		Market: market,
		CalibratedSlice: &pb.ImpliedVolSlice{
			Timestamp:         market.Timestamp,
			Expiry:            quotes.Expiry,
			Iserror:           false,
			Vols:              make([]float64, n),
			Strikes:           make([]float64, n),
			Errors:            make([]string, n),
			Quotes:            make([]*pb.OptionQuote, n),
			Nbsolveriteration: make([]int64, n),
		},
	}
}

func (it optionQuoteSliceIterator) update(quote *pb.OptionQuote, res *ivSolverResult) {
	if res.Error == nil {
		it.CalibratedSlice.Quotes[it.NbCalibratedSlice] = quote
		it.CalibratedSlice.Vols[it.NbCalibratedSlice] = res.IV
		it.CalibratedSlice.Strikes[it.NbCalibratedSlice] = quote.Strike
		it.CalibratedSlice.Nbsolveriteration[it.NbCalibratedSlice] = int64(res.NbSolverIteration)
	} else {
		it.CalibratedSlice.Iserror = true
		it.CalibratedSlice.Errors[it.NbCalibratedSlice] = res.Error.Error()
	}

	it.NbCalibratedSlice++
}

func (it optionQuoteSliceIterator) foreach(f func(quote *pb.OptionQuote) *ivSolverResult) *optionQuoteSliceIterator {
	it.NbCalibratedSlice = 0
	spot := it.Market.Spot.Index.Value

	for index := 0; index < len(it.Slice.Puts); index++ {
		if index < len(it.Slice.Puts) && it.Slice.Puts[index].Strike <= spot {
			if it.Slice.Puts[index].Strike/spot > putLBound {
				res := f(it.Slice.Puts[index])
				it.update(it.Slice.Puts[index], res)
			}
		}
	}

	for index := 0; index < len(it.Slice.Calls); index++ {
		if index < len(it.Slice.Calls) && it.Slice.Calls[index].Strike > spot {
			res := f(it.Slice.Calls[index])
			it.update(it.Slice.Calls[index], res)
		}
	}

	it.CalibratedSlice.Vols = it.CalibratedSlice.Vols[0:it.NbCalibratedSlice]
	it.CalibratedSlice.Nbsolveriteration = it.CalibratedSlice.Nbsolveriteration[0:it.NbCalibratedSlice]
	it.CalibratedSlice.Strikes = it.CalibratedSlice.Strikes[0:it.NbCalibratedSlice]
	it.CalibratedSlice.Quotes = it.CalibratedSlice.Quotes[0:it.NbCalibratedSlice]
	it.CalibratedSlice.Errors = it.CalibratedSlice.Errors[0:it.NbCalibratedSlice]
	return &it
}

func (it optionQuoteSliceIterator) then(f func(calibratedSlice *pb.ImpliedVolSlice)) *optionQuoteSliceIterator {
	f(it.CalibratedSlice)
	return &it
}
