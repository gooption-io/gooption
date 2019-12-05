package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"github.com/lehajam/dgo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "github.com/gooption-io/gooption/v1/contract/api"
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
	query := `{
		Europeans(func: eq(type, "european")){
			uid
			ticker
			strike
			undticker
			expiry
			putcall
		}
	}`

	resp, err := s.db.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	match := &api_pb.ListEuropeansResponse{}
	err = dgo.Unmarshal(resp.GetJson(), match)
	if err != nil {
		return nil, err
	}

	return match, nil
}

func (s *europeanServiceServerImpl) GetEuropean(ctx context.Context, req *api_pb.GetEuropeanRequest) (*api_pb.European, error) {
	query := `query GetEuropean($optionTicker: string){
		Europeans(func: eq(ticker, $optionTicker)){
			uid
			ticker
			strike
			undticker
			expiry
			putcall
		}
	}`

	params := map[string]string{
		"$optionTicker": req.Ticker,
	}

	resp, err := s.db.NewTxn().QueryWithVars(ctx, query, params)
	if err != nil {
		return nil, err
	}

	match := &api_pb.ListEuropeansResponse{}
	err = dgo.Unmarshal(resp.GetJson(), match)
	if err != nil {
		return nil, err
	}

	if len(match.Europeans) == 0 {
		return nil, nil
	}

	return match.Europeans[0], nil
}

type typeDecorator struct {
	UID  string `json:"uid"`
	Type string `json:"type"`
}

func (s *europeanServiceServerImpl) CreateEuropean(ctx context.Context, req *api_pb.CreateEuropeanRequest) (*api_pb.European, error) {
	assigned, err := insertObj(ctx, s.db, req.European)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	withType := typeDecorator{assigned.Uids["blank-0"], "european"}
	_, err = insertObj(ctx, s.db, withType)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	req.European.Uid = withType.UID
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

func createSchema(ctx context.Context, db *dgo.Dgraph) error {
	err := alterSchema(ctx, db,
		`timestamp: float @index(float) .
		ticker: string @index(exact, term) .
		undticker: string @index(exact, term) @upsert .
		type: string @index(exact, term) .`)

	if err != nil {
		logrus.Errorln(err)
		return err
	}

	return nil
}
