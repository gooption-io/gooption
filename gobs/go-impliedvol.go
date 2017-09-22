package gooption

import (
	"errors"
	"sort"
	"time"

	"github.com/gooption/pb"
)

type VolSlice struct {
	Expiry          float64
	Strikes, Values []float64
	OpenInterest    []float64
	Errors          []string
	IsError         bool
}

func NewVolSlice(expiry float64, NbContract int) VolSlice {
	return VolSlice{
		Expiry:       expiry,
		Strikes:      make([]float64, NbContract),
		Values:       make([]float64, NbContract),
		OpenInterest: make([]float64, NbContract),
		Errors:       make([]string, NbContract),
		IsError:      false,
	}
}

func buildVolSurface(pricingDate time.Time, spot, rate float64, optionQuotes []*pb.ImpliedVolRequest_QuoteTermStructure) ([]VolSlice, error) {
	volSurface := make([]VolSlice, len(optionQuotes))
	for i, quote := range optionQuotes {
		life := time.Unix(int64(quote.Expiry), 0).Sub(pricingDate).Hours() / 24.0 / 365.25
		slice, err := buildSlice(
			spot,
			rate,
			life,
			quote.Calls,
			quote.Puts)
		if err != nil {
			slice = VolSlice{IsError: true, Errors: []string{err.Error()}}
		}

		volSurface[i] = slice
	}
	return volSurface, nil
}

func buildSlice(
	spot, rate, life float64,
	calls []*pb.ImpliedVolRequest_Quote,
	puts []*pb.ImpliedVolRequest_Quote) (VolSlice, error) {

	if len(calls) < 3 || len(puts) < 3 {
		return VolSlice{}, errors.New("Not enough quotes to calibrate")
	}

	callLowerBound := sort.Search(len(calls), func(i int) bool { return calls[i].Strike >= spot })
	if callLowerBound < 2 {
		return VolSlice{}, errors.New("Not enough Call quotes to calibrate")
	}

	putUpperBound := sort.Search(len(puts), func(i int) bool { return puts[i].Strike >= calls[callLowerBound].Strike })
	if putUpperBound < 2 {
		return VolSlice{}, errors.New("Not enough Put quotes to calibrate")
	}

	putUpperBound = putUpperBound - 1
	putLowerBound := sort.Search(putUpperBound, func(i int) bool { return puts[i].Strike/spot >= 0.20 })
	if putLowerBound >= putUpperBound {
		return VolSlice{}, errors.New("Put moneyness must be greater than 20%")
	}

	slice := NewVolSlice(life, len(calls)-callLowerBound+putUpperBound-putLowerBound)
	calibrateSlice(putLowerBound, putUpperBound, 0, spot, rate, life, -1.0, &slice, puts)
	calibrateSlice(callLowerBound, len(calls), putUpperBound-putLowerBound, spot, rate, life, 1.0, &slice, calls)

	return slice, nil
}

func calibrateSlice(
	lbound, ubound, sliceIndex int,
	spot, rate, life, putCall float64,
	slice *VolSlice,
	quotes []*pb.ImpliedVolRequest_Quote) {

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
