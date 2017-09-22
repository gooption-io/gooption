package marketdata

import (
	"time"

	"github.com/gooption/pb"
)

type Service struct {
	Market *pb.ImpliedVolMarket, 
	Quotes *OptionQuoteSurface,
}

func NewImpliedVolService(market *pb.ImpliedVolMarket, quotes *OptionQuoteSurface) {
	return {
		Market: market,
		Quotes: quotes
	}
}

func (s Service) Calibrate() (*pb.ImpliedVolSurface, error) {
	volSurface, err := NewImpliedVolSurface(Market, Quotes)
	if err != nil {
		return nil, err
	}
	return volSurface, nil

	// ivsurface = make([]*pb.ImpliedVolSlice, len(volSurface)),
	// for i, slice := range volSurface {
	// 	ivsurface[i] = &pb.ImpliedVolSlice{
	// 		Iserror: slice.IsError,
	// 		Expiry:  slice.Expiry,
	// 		Strikes: slice.Strikes,
	// 		Vols:    slice.Values,
	// 		Errors:  slice.Errors,
	// 	}
	// }
}

func NewImpliedVolSurface(market *pb.ImpliedVolMarket, quotes *OptionQuoteSurface) ([]VolSlice, error) {	
	pricingDate = time.Unix(Market.timestamp)
	volSurface := make([]VolSlice, len(optionQuotes))
	for i, quoteSlice := range slices {
		slice, err := NewCalibratedSlice(pricingDate, market, quoteSlice)
		if err != nil {
			slice = VolSlice{IsError: true, Errors: []string{err.Error()}}
		}

		volSurface[i] = slice
	}
	return volSurface, nil
}

type PutCallBoundary struct {
	LBound, UBound int
}

func NewPutCallBoundary(quoteSlice *OptionQuoteSlice) (PutCallBoundary, PutCallBoundary) {
	dataFilter := func(i int) bool { 
		return quoteSlice.calls[i].Strike >= spot 
	}

	callFilter := sort.Search(len(putCallQuotes), func(i int) bool { 
		return quoteSlice.calls[i].Strike >= spot 
	})

	callBoundary := sort.Search(len(putCallQuotes), func(i int) bool { 
		return quoteSlice.puts[i].Strike >= quoteSlice.calls[callLowerBound].Strike 
		&& quoteSlice.puts[i].Strike/spot >= 0.20 
		})

	if putBoundary < 2 || callBoundary < 2 {
		return VolSlice{}, errors.New("Not enough " + putBoundary < 2 ? "puts" : "calls" + " quotes to calibrate")
	}

	return PutCallBoundary { LBound: putBoundary, UBound: } (putBoundary, callBoundary)
}


func NewCalibratedSlice(pricingDate time.Time, market *pb.ImpliedVolMarket, quoteSlice *OptionQuoteSlice) (VolSlice, error) {
	if len(calls) < 3 || len(puts) < 3 {
		return VolSlice{}, errors.New("Not enough quotes to calibrate")
	}

	life := time.Unix(int64(quote.Expiry), 0).Sub(pricingDate).Hours() / 24.0 / 365.25
	callBoundary, putBoundary := NewPutCallBoundary(quoteSlice);
	slice := NewVolSlice(life, len(calls)-callLowerBound+putUpperBound-putLowerBound)
	NewCalibratedSlice(slice[putLowerBound:putBoundary], 0, spot, rate, life, -1.0, &slice, puts)
	NewCalibratedSlice(slice[callLowerBound:len(calls)], putUpperBound-putLowerBound, spot, rate, life, 1.0, &slice, calls)

	return slice, nil
}

func NewQuoteIterator(quoteSlice *OptionQuoteSlice, ) func() *OptionQuote {
	for index := lbound; index < ubound; index++ {
		iv, _, err := ImpliedVol(
			quotes[index].Bid,
			spot,
			quotes[index].Strike,
			life,
			rate,
			putCall)
		if err != nil {
			slice.IsError = true
			slice.Errors[currentSliceIndex] = err.Error()
		} else {
			slice.Strikes[currentSliceIndex] = quotes[index].Strike
			slice.Values[currentSliceIndex] = iv
		}
		currentSliceIndex++
	}	
}

func NewCalibratorIterator() {
	currentSliceIndex := sliceIndex
	for index := lbound; index < ubound; index++ {
		iv, _, err := ImpliedVol(
			quotes[index].Bid,
			spot,
			quotes[index].Strike,
			life,
			rate,
			putCall)
		if err != nil {
			slice.IsError = true
			slice.Errors[currentSliceIndex] = err.Error()
		} else {
			slice.Strikes[currentSliceIndex] = quotes[index].Strike
			slice.Values[currentSliceIndex] = iv
		}
		currentSliceIndex++
	}	
}

func calibrateSlice(pricingDate time.Time, market *pb.ImpliedVolMarket, quoteSlice *OptionQuoteSlice)  {

	currentSliceIndex := sliceIndex
	quotes := NewQuote(quoteSlice)
	while(quote := quotes()) {

	}

	for index := lbound; index < ubound; index++ {
		iv, _, err := ImpliedVol(
			quotes[index].Bid,
			spot,
			quotes[index].Strike,
			life,
			rate,
			putCall)
		if err != nil {
			slice.IsError = true
			slice.Errors[currentSliceIndex] = err.Error()
		} else {
			slice.Strikes[currentSliceIndex] = quotes[index].Strike
			slice.Values[currentSliceIndex] = iv
		}
		currentSliceIndex++
	}
}
