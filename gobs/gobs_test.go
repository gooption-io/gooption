package main

import (
	"math"
	"testing"
	"time"

	"github.com/gooption-io/gooption/proto/go/pb"
)

var (
	pricingDate float64 = 1540425600

	call = pb.European{
		Strike:  100,
		Putcall: "call",
		Expiry:  1571961600,
	}

	put = pb.European{
		Strike:  100,
		Putcall: "put",
		Expiry:  1571961600,
	}

	mkt = pb.OptionMarket{
		Timestamp: 1540425600,
		Spot:      &pb.Spot{&pb.Index{Value: 100}},
		Vol:       &pb.FlatVol{&pb.Index{Value: 0.1}},
		Rate:      &pb.RiskFreeRate{&pb.Index{Value: 0.01}},
	}
)

func BenchmarkFairValue(b *testing.B) {
	for index := 0; index < b.N; index++ {
		bs(pricingDate, &call, &mkt)
	}
}

func TestPutCallParity(t *testing.T) {
	var (
		s = mkt.Spot.Index.Value
		r = mkt.Rate.Index.Value
		k = call.Strike
		T = time.Unix(int64(call.Expiry), 0).Sub(
			time.Unix(int64(pricingDate), 0)).Hours() / 24.0 / 365.250
	)

	c := bs(pricingDate, &call, &mkt)
	p := bs(pricingDate, &put, &mkt)
	f := s - math.Exp(-r*T)*k
	t.Logf("Forward: %v", f)

	if math.Abs(c-p-f) > 1E-10 {
		t.Errorf("Put Call parity broken")
		t.Logf("Call - Put: %v", c-p)
	}
}

func TestPutCallIVRootSolver(t *testing.T) {
	c := bs(pricingDate, &call, &mkt)
	t.Logf("atmCall price: %v", c)

	quote := pb.OptionQuote{
		Ask:     c,
		Bid:     c,
		Strike:  call.Strike,
		Putcall: call.Putcall,
	}

	iv := ivRootSolver(
		pricingDate,
		call.Expiry,
		&quote,
		&mkt,
	)
	if iv.Error != "" {
		t.Errorf(iv.Error)
	}

	t.Logf("atmCall iv: %v", iv.Vol)
	t.Logf("atmCall iv iteration: %v", iv.Nbiteration)
	if math.Abs(iv.Vol-mkt.Vol.Index.Value) > 1E-10 {
		t.Errorf("atmCall iv %v should be equal to %v", iv.Vol, mkt.Vol.Index.Value)
	}

	p := bs(pricingDate, &put, &mkt)
	t.Logf("atmPut price: %v", p)

	quote = pb.OptionQuote{
		Ask:     p,
		Bid:     p,
		Strike:  put.Strike,
		Putcall: put.Putcall,
	}

	iv = ivRootSolver(
		pricingDate,
		put.Expiry,
		&quote,
		&mkt,
	)
	if iv.Error != "" {
		t.Errorf(iv.Error)
	}

	t.Logf("atmPut iv: %v", iv.Vol)
	t.Logf("atmPut iv iteration: %v", iv.Nbiteration)
	if math.Abs(iv.Vol-mkt.Vol.Index.Value) > 1E-10 {
		t.Errorf("atmPut iv %v should be equal to %v", iv.Vol, mkt.Vol.Index.Value)
	}
}
