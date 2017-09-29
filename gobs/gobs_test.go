package main

import (
	"math"
	"testing"
)

var (
	params = map[string]float64{
		"S":       100,
		"K":       100,
		"T":       1,
		"Sigma":   0.1,
		"R":       0.01,
		"Call":    1.0,
		"Put":     -1.0,
		"atmCall": 4.485236409022086,
		"atmPut":  3.4902197839388975,
	}
)

func BenchmarkFairValue(b *testing.B) {
	for index := 0; index < b.N; index++ {
		bs(
			params["S"],
			params["Sigma"],
			params["R"],
			params["K"],
			params["T"],
			params["Call"])
	}
}
func TestPutCallParity(t *testing.T) {
	call := bs(
		params["S"],
		params["Sigma"],
		params["R"],
		params["K"],
		params["T"],
		params["Call"])
	t.Logf("Call: %v", call)

	put := bs(
		params["S"],
		params["Sigma"],
		params["R"],
		params["K"],
		params["T"],
		params["Put"])
	t.Logf("Put: %v", put)

	forward := params["S"] - math.Exp(-params["R"]*params["T"])*params["K"]
	t.Logf("Forward: %v", forward)

	if math.Abs(call-put-forward) > 1E-10 {
		t.Errorf("Put Call parity broken")
		t.Logf("Call - Put: %v", call-put)
	}
}

func TestPutCallIVRootSolver(t *testing.T) {
	call := ivRootSolver(
		params["atmCall"],
		params["S"],
		params["R"],
		params["K"],
		params["T"],
		params["Call"])
	t.Logf("atmCall iv: %v", call.IV)
	t.Logf("atmCall iv iteration: %v", call.NbSolverIteration)
	if call.Error != nil {
		t.Errorf(call.Error.Error())
	}
	if math.Abs(call.IV-params["Sigma"]) > 1E-10 {
		t.Errorf("atmCall iv %v should be equal to %v", call.IV, params["Sigma"])
	}

	put := ivRootSolver(
		params["atmPut"],
		params["S"],
		params["R"],
		params["K"],
		params["T"],
		params["Put"])
	t.Logf("atmPut iv: %v", put.IV)
	t.Logf("atmPut iv iteration: %v", put.NbSolverIteration)
	if put.Error != nil {
		t.Errorf(put.Error.Error())
	}
	if math.Abs(put.IV-params["Sigma"]) > 1E-10 {
		t.Errorf("atmPut iv %v should be equal to %v", put.IV, params["Sigma"])
	}

	if math.Abs(call.IV-put.IV) > 1E-10 {
		t.Errorf("atm vol call different from atm put")
	}
}
