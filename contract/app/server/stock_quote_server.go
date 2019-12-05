package server

import (
	"context"

	"github.com/lehajam/dgo"
	"github.com/sirupsen/logrus"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "github.com/gooption-io/gooption/v1/contract/api"
)

// StockQuoteServiceServer is a composite interface of api_pb.StockQuoteServiceServer and grapiserver.Server.
type StockQuoteServiceServer interface {
	api_pb.StockQuoteServiceServer
	grapiserver.Server
}

// NewStockQuoteServiceServer creates a new StockQuoteServiceServer instance.
func NewStockQuoteServiceServer() StockQuoteServiceServer {
	return &stockQuoteServiceServerImpl{
		db: newDgraphClient(),
	}
}

type stockQuoteServiceServerImpl struct {
	db *dgo.Dgraph
}

func (s *stockQuoteServiceServerImpl) ListStockQuotes(ctx context.Context, req *api_pb.ListStockQuotesRequest) (*api_pb.ListStockQuotesResponse, error) {
	query := `{
	  StockQuotes(func: eq(type, "stock_quote")){
		uid
		timestamp
		ticker
		open
		close
		low
		high
		volume
	  }
	}`

	resp, err := s.db.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	match := &api_pb.ListStockQuotesResponse{}
	err = dgo.Unmarshal(resp.GetJson(), match)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (s *stockQuoteServiceServerImpl) GetStockQuote(ctx context.Context, req *api_pb.GetStockQuoteRequest) (*api_pb.StockQuote, error) {
	query := `query GetStockQuote($uid: string){
		StockQuotes(func: uid($uid)){
			uid
			timestamp
			ticker
			open
			close
			low
			high
			volume
		}
	}`

	params := map[string]string{
		"$uid": req.StockQuoteId,
	}

	resp, err := s.db.NewTxn().QueryWithVars(ctx, query, params)
	if err != nil {
		return nil, err
	}

	match := &api_pb.ListStockQuotesResponse{}
	err = dgo.Unmarshal(resp.GetJson(), match)
	if err != nil {
		return nil, err
	}

	if len(match.StockQuotes) == 0 {
		return nil, nil
	}

	return match.StockQuotes[0], nil
}

func (s *stockQuoteServiceServerImpl) CreateStockQuote(ctx context.Context, req *api_pb.CreateStockQuoteRequest) (*api_pb.StockQuote, error) {
	assigned, err := insertObj(ctx, s.db, req.StockQuote)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	withType := typeDecorator{assigned.Uids["blank-0"], "stock_quote"}
	_, err = insertObj(ctx, s.db, withType)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	req.StockQuote.Uid = withType.UID
	return req.StockQuote, nil
}

func (s *stockQuoteServiceServerImpl) UpdateStockQuote(ctx context.Context, req *api_pb.UpdateStockQuoteRequest) (*api_pb.StockQuote, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *stockQuoteServiceServerImpl) DeleteStockQuote(ctx context.Context, req *api_pb.DeleteStockQuoteRequest) (*empty.Empty, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
