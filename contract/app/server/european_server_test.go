package server

import (
	"context"
	"testing"
	"time"

	api_pb "github.com/gooption-io/gooption/v1/contract/api"

	"github.com/stretchr/testify/require"
)

func Test_EuropeanServiceServer_ListEuropeans(t *testing.T) {
	svr := NewEuropeanServiceServer()

	ctx := context.Background()
	req := &api_pb.ListEuropeansRequest{}

	resp, err := svr.ListEuropeans(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Europeans)
}

func Test_EuropeanServiceServer_GetEuropean(t *testing.T) {
	svr := NewEuropeanServiceServer()

	ctx := context.Background()
	req := &api_pb.GetEuropeanRequest{
		Ticker: "AAPL DEC2019 PUT",
	}

	resp, err := svr.GetEuropean(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Uid)
}

func Test_EuropeanServiceServer_CreateEuropean(t *testing.T) {
	svr := NewEuropeanServiceServer()

	ctx := context.Background()

	req := &api_pb.CreateEuropeanRequest{
		European: &api_pb.European{
			Timestamp: 1514162664,
			Ticker:    "AAPL DEC2019 PUT",
			Undticker: "AAPL",
			Strike:    159.76,
			Expiry:    float64(time.Now().AddDate(0, 1, 0).Unix()),
			Putcall:   "put",
		},
	}

	resp, err := svr.CreateEuropean(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Uid)
}

func Test_EuropeanServiceServer_UpdateEuropean(t *testing.T) {
	svr := NewEuropeanServiceServer()

	ctx := context.Background()
	req := &api_pb.UpdateEuropeanRequest{}

	resp, err := svr.UpdateEuropean(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_EuropeanServiceServer_DeleteEuropean(t *testing.T) {
	svr := NewEuropeanServiceServer()

	ctx := context.Background()
	req := &api_pb.DeleteEuropeanRequest{}

	resp, err := svr.DeleteEuropean(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_europeanServiceServerImpl_CreateSchema(t *testing.T) {
	ctx := context.Background()
	if err := createSchema(ctx, newDgraphClient()); err != nil {
		t.Errorf("europeanServiceServerImpl.CreateSchema() error = %v", err)
	}
}
