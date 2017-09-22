package main

import (
	"testing"

	"github.com/goyahoo"
)

func Test_buildSurface(t *testing.T) {
	ticker := "AAPL"
	chain, _, _ := goyahoo.GetOptionChain(ticker)
	volSurface := buildVolSurface(chain)
	drawSurface(ticker, volSurface)
}
