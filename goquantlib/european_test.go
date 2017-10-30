package gooption

import (
	"testing"
)

var (
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
	call := EuropeanFlatVol(
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
