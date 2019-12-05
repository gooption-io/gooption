package server

import (
	api_pb "gobs/api"
	"testing"
	"context"
)

func Test_europeanPricerServiceServerImpl_Compute(t *testing.T) {
	svr := NewEuropeanPricerServiceServer()
	ctx := context.Background()

	svr.Compute(ctx, &api_pb.ComputationRequest{})
}