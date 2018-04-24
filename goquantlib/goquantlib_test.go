package main

import (
	"context"
	"os"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/lehajam/gooption/goquantlib/pb"
	"github.com/lehajam/gooption/goquantlib/quantlib"
)

var (
	s      = &server{}
	params = map[string]float64{
		"S":       100,
		"K":       100,
		"Sigma":   0.1,
		"R":       0.01,
		"Q":       0,
		"T0":      1508705500,
		"T":       1524182400,
		"Call":    1.0,
		"Put":     -1.0,
		"atmCall": 4.485236409022086,
		"atmPut":  3.4902197839388975,
	}
)

func TestPrice(t *testing.T) {
	call := quantlib.EuropeanFlatVol(
		params["S"],
		params["R"],
		params["Q"],
		params["Sigma"],
		params["K"],
		int(params["T0"]),
		int(params["T"]),
		int(params["Call"]))
	t.Logf("Call: %v", call)
}

func Test_Price(t *testing.T) {
	if file, err := os.Open("./testdata/PriceRequest.json"); err == nil {
		defer file.Close()
		request := &pb.PriceRequest{}
		if jsonpb.Unmarshal(file, request) == nil {
			if response, err := s.Price(context.Background(), request); err != nil {
				t.Log(response)
				t.Error(err)
				if response.Error != "" {
					t.Error(response.Error)
				}
			}
		} else {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}
}
