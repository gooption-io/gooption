package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"github.com/lehajam/dgo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "contract/api"
)

// EuropeanServiceServer is a composite interface of api_pb.EuropeanServiceServer and grapiserver.Server.
type EuropeanServiceServer interface {
	api_pb.EuropeanServiceServer
	grapiserver.Server
}

// NewEuropeanServiceServer creates a new EuropeanServiceServer instance.
func NewEuropeanServiceServer() EuropeanServiceServer {
	return &europeanServiceServerImpl{
		db: newDgraphClient(),
	}
}

type europeanServiceServerImpl struct {
	db *dgo.Dgraph
}

func (s *europeanServiceServerImpl) ListEuropeans(ctx context.Context, req *api_pb.ListEuropeansRequest) (*api_pb.ListEuropeansResponse, error) {
	contracts := getAllContracts(ctx, s.db)
	return &api_pb.ListEuropeansResponse{
		Europeans: contracts,
	}, nil
}

func (s *europeanServiceServerImpl) GetEuropean(ctx context.Context, req *api_pb.GetEuropeanRequest) (*api_pb.European, error) {
	resp := query(ctx, s.db,
		`query GetEuropean($optionTicker: string){
			contract(func: eq(ticker, $optionTicker)){
				ticker
				strike
				und as undticker
				expiry
				putcall
			}
		}`,
		map[string]string{
			"$optionTicker": req.Ticker,
		})

	european := &api_pb.European{}
	err := dgo.Unmarshal(resp.GetJson(), european)
	if err != nil {
		panic(err)
	}

	// TODO: Not yet implemented.
	return european, nil
}

func (s *europeanServiceServerImpl) CreateEuropean(ctx context.Context, req *api_pb.CreateEuropeanRequest) (*api_pb.European, error) {
	insert(ctx, s.db, req.European,
		`timestamp: float @index(float) .
		ticker: string @index(exact, term) .
		undticker: string @index(exact, term) .`,
	)

	return req.European, nil
}

func (s *europeanServiceServerImpl) UpdateEuropean(ctx context.Context, req *api_pb.UpdateEuropeanRequest) (*api_pb.European, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *europeanServiceServerImpl) DeleteEuropean(ctx context.Context, req *api_pb.DeleteEuropeanRequest) (*empty.Empty, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
