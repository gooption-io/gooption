package server

import (
	"context"
	"testing"

	api_pb "github.com/gooption-io/gooption/api"
)

func Test_StockServiceServer_ListStocks(t *testing.T) {
	svr := NewStockServiceServer()

	ctx := context.Background()
	req := &api_pb.ListStocksRequest{}

	resp, err := svr.ListStocks(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_StockServiceServer_GetStock(t *testing.T) {
	svr := NewStockServiceServer()

	ctx := context.Background()
	req := &api_pb.GetStockRequest{}

	resp, err := svr.GetStock(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_StockServiceServer_CreateStock(t *testing.T) {
	svr := NewStockServiceServer()

	ctx := context.Background()
	req := &api_pb.CreateStockRequest{}

	resp, err := svr.CreateStock(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_StockServiceServer_UpdateStock(t *testing.T) {
	svr := NewStockServiceServer()

	ctx := context.Background()
	req := &api_pb.UpdateStockRequest{}

	resp, err := svr.UpdateStock(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_StockServiceServer_DeleteStock(t *testing.T) {
	svr := NewStockServiceServer()

	ctx := context.Background()
	req := &api_pb.DeleteStockRequest{}

	resp, err := svr.DeleteStock(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}
