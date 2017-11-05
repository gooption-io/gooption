package main

import (
	"sync/atomic"

	"github.com/gooption/gobs/pb"
)

var (
	putLBound = 0.20
)

type optionQuoteSliceIterator struct {
	NbCalibratedSlice int32
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
			Timestamp: market.Timestamp,
			Expiry:    quotes.Expiry,
			Iserror:   false,
			Quotes:    make([]*pb.ImpliedVolQuote, n),
		},
	}
}

func (it *optionQuoteSliceIterator) update(quote *pb.OptionQuote, res *ivSolverResult) {
	if res.Error == nil {
		it.CalibratedSlice.Quotes[it.NbCalibratedSlice] = &pb.ImpliedVolQuote{
			Input:       quote,
			Vol:         res.IV,
			Nbiteration: int64(res.NbSolverIteration),
		}
	} else {
		it.CalibratedSlice.Iserror = true
		it.CalibratedSlice.Quotes[it.NbCalibratedSlice] = &pb.ImpliedVolQuote{
			Input: quote,
			Error: res.Error.Error(),
		}
	}

	atomic.AddInt32(&it.NbCalibratedSlice, 1)
}

func (it *optionQuoteSliceIterator) foreach(f func(quote *pb.OptionQuote) *ivSolverResult) *optionQuoteSliceIterator {
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

	it.CalibratedSlice.Quotes = it.CalibratedSlice.Quotes[0:it.NbCalibratedSlice]
	return it
}

func (it *optionQuoteSliceIterator) then(f func(calibratedSlice *pb.ImpliedVolSlice)) *optionQuoteSliceIterator {
	f(it.CalibratedSlice)
	return it
}
