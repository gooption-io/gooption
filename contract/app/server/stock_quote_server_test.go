package server

import (
	"context"
	"testing"

	api_pb "contract/api"
)

func Test_StockQuoteServiceServer_ListStockQuotes(t *testing.T) {
	svr := NewStockQuoteServiceServer()

	ctx := context.Background()
	req := &api_pb.ListStockQuotesRequest{}

	resp, err := svr.ListStockQuotes(ctx, req)

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}

	if len(resp.StockQuotes) == 0 {
		t.Error("no data found")
	}
}

func Test_StockQuoteServiceServer_GetStockQuote(t *testing.T) {
	svr := NewStockQuoteServiceServer()

	ctx := context.Background()
	req := &api_pb.GetStockQuoteRequest{
		StockQuoteId: "0x80ac",
	}

	resp, err := svr.GetStockQuote(ctx, req)

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_StockQuoteServiceServer_CreateStockQuote(t *testing.T) {
	svr := NewStockQuoteServiceServer()

	ctx := context.Background()
	req := &api_pb.CreateStockQuoteRequest{
		StockQuote: &api_pb.StockQuote{
			Timestamp: 1575244800000,
			Ticker:    "AAPL",
			Open:      267.27,
			Close:     264.16,
			Volume:    22461876,
			Low:       263.45,
			High:      268.25,
		},
	}

	resp, err := svr.CreateStockQuote(ctx, req)

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}

	if resp.Uid == "" {
		t.Error("response UID should not nil")
	}
}

func Test_StockQuoteServiceServer_UpdateStockQuote(t *testing.T) {
	svr := NewStockQuoteServiceServer()

	ctx := context.Background()
	req := &api_pb.UpdateStockQuoteRequest{}

	resp, err := svr.UpdateStockQuote(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_StockQuoteServiceServer_DeleteStockQuote(t *testing.T) {
	svr := NewStockQuoteServiceServer()

	ctx := context.Background()
	req := &api_pb.DeleteStockQuoteRequest{}

	resp, err := svr.DeleteStockQuote(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}
