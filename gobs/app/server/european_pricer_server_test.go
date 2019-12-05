package server

import (
	"context"
	"testing"

	api_pb "github.com/gooption-io/gooption/v1/gobs/api"
)

func Test_europeanPricerServiceServerImpl_Compute(t *testing.T) {
	svr := NewEuropeanPricerServiceServer()
	ctx := context.Background()

	svr.Compute(ctx, &api_pb.ComputationRequest{})
}
