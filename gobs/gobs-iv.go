package gooption

import (
	"github.com/gooption/pb"
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

func (it optionQuoteSliceIterator) update(quote *pb.OptionQuote, res *IVSolverResult) {
	if res.Error == nil {
		it.CalibratedSlice.Quotes[it.NbCalibratedSlice] = quote
		it.CalibratedSlice.Vols[it.NbCalibratedSlice] = res.IV
		it.CalibratedSlice.Strikes[it.NbCalibratedSlice] = quote.Strike
		it.CalibratedSlice.Nbsolveriteration[it.NbCalibratedSlice] = int64(res.NbSolverIteration)
	} else {
		it.CalibratedSlice.Iserror = true
		it.CalibratedSlice.Errors[it.NbCalibratedSlice] = res.Error.Error()
	}
}

func (it optionQuoteSliceIterator) foreach(f func(quote *pb.OptionQuote) *IVSolverResult) *optionQuoteSliceIterator {
	it.NbCalibratedSlice = 0
	increment := func(q *pb.OptionQuote) {
		res := f(q)
		it.update(q, res)
		it.NbCalibratedSlice++
	}

	for index := 0; index < len(it.Slice.Puts); index++ {
		if index < len(it.Slice.Puts) && it.Slice.Puts[index].Strike <= it.Market.Spot.Value {
			if it.Slice.Puts[index].Strike/it.Market.Spot.Value > 0.20 {
				increment(it.Slice.Puts[index])
			}
		}
	}

	for index := 0; index < len(it.Slice.Calls); index++ {
		if index < len(it.Slice.Calls) && it.Slice.Calls[index].Strike > it.Market.Spot.Value {
			increment(it.Slice.Calls[index])
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
