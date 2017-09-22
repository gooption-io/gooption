package main

import (
	"bytes"
	"sort"

	"time"

	"io/ioutil"

	"errors"

	"github.com/gooption"
	"github.com/goyahoo"
	"github.com/mkideal/cli"
	chart "github.com/wcharczuk/go-chart"
)

type argT struct {
	cli.Helper
	Ticker string `cli:"t,ticker" usage:"ticker for the option underlier"`
	// PricingDate string `cli:"d, pricingDate" usage:"pricing date"`
}

type VolSlice struct {
	Expiry          float64
	Strikes, Values []float64
	OpenInterest    []float64
	Errors          []string
	IsError         bool
}

func NewVolSlice(NbContract int) VolSlice {
	return VolSlice{
		Strikes:      make([]float64, NbContract),
		Values:       make([]float64, NbContract),
		OpenInterest: make([]float64, NbContract),
		Errors:       make([]string, NbContract),
		IsError:      false,
	}
}

func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)
		ctx.String("ticker=%s\n", argv.Ticker)

		chain, _, _ := goyahoo.GetOptionChain(argv.Ticker)
		volSurface := buildVolSurface(chain)
		drawSurface(argv.Ticker, volSurface)

		return nil
	})
}

func buildVolSurface(OptionChain []*goyahoo.Chain) []VolSlice {
	volSurface := make([]VolSlice, len(OptionChain))
	for i, chain := range OptionChain {
		life := time.Unix(chain.ExpirationDates[i], 0).Sub(time.Now()).Hours() / 24.0 / 365.25
		volSurface[i], _ = buildSlice(
			chain.Quote.RegularMarketPrice,
			life,
			chain.Options[0].Calls,
			chain.Options[0].Puts)
	}
	return volSurface
}

func buildSlice(
	spot, life float64,
	calls []goyahoo.Quote,
	puts []goyahoo.Quote) (VolSlice, error) {

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

	slice := NewVolSlice(len(calls) - callLowerBound + putUpperBound - putLowerBound)
	calibrateSlice(putLowerBound, putUpperBound, 0, spot, life, -1.0, &slice, puts)
	calibrateSlice(callLowerBound, len(calls), putUpperBound-putLowerBound, spot, life, 1.0, &slice, calls)

	return slice, nil
}

func calibrateSlice(
	lbound, ubound, sliceIndex int,
	spot, life, putCall float64,
	slice *VolSlice, quotes []goyahoo.Quote) {

	currentSliceIndex := sliceIndex
	for index := lbound; index < ubound; index++ {
		iv, _, err := gooption.ImpliedVol(
			quotes[index].LastPrice,
			spot,
			quotes[index].Strike,
			life,
			0.01,
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

func drawSurface(Ticker string, Slices []VolSlice) {
	nbSlices := 0
	for _, slice := range Slices {
		if !slice.IsError {
			nbSlices++
		}
	}

	graph := chart.Chart{
		Series: make([]chart.Series, nbSlices),
	}

	sliceIndex := 0
	for _, slice := range Slices {
		if !slice.IsError {
			graph.Series[sliceIndex] = chart.ContinuousSeries{
				XValues: slice.Strikes,
				YValues: slice.Values,
			}
			sliceIndex++
		}
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)
	ioutil.WriteFile(Ticker+".png", buffer.Bytes(), 0644)
}

func drawSlice(Ticker string, Slice VolSlice) {
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: Slice.Strikes,
				YValues: Slice.Values,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)
	ioutil.WriteFile(Ticker+".png", buffer.Bytes(), 0644)
}
