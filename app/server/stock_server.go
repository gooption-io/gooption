package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "github.com/gooption-io/gooption/api"
)

// StockServiceServer is a composite interface of api_pb.StockServiceServer and grapiserver.Server.
type StockServiceServer interface {
	api_pb.StockServiceServer
	grapiserver.Server
}

// NewStockServiceServer creates a new StockServiceServer instance.
func NewStockServiceServer() StockServiceServer {
	return &stockServiceServerImpl{}
}

type stockServiceServerImpl struct {
}

func (s *stockServiceServerImpl) ListStocks(ctx context.Context, req *api_pb.ListStocksRequest) (*api_pb.ListStocksResponse, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *stockServiceServerImpl) GetStock(ctx context.Context, req *api_pb.GetStockRequest) (*api_pb.Stock, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *stockServiceServerImpl) CreateStock(ctx context.Context, req *api_pb.CreateStockRequest) (*api_pb.Stock, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *stockServiceServerImpl) UpdateStock(ctx context.Context, req *api_pb.UpdateStockRequest) (*api_pb.Stock, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *stockServiceServerImpl) DeleteStock(ctx context.Context, req *api_pb.DeleteStockRequest) (*empty.Empty, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
