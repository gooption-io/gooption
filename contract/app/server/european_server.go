package server

import (
	"context"
	"fmt"

	dgo "github.com/dgraph-io/dgo/v2"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/izumin5210/grapi/pkg/grapiserver"
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

// used to allow * queries by types
// to be replaced by dgraph types
type typeDecorator struct {
	UID  string `json:"uid"`
	Type string `json:"type"`
}

type ListEuropeansResponseVals struct {
	Europeans []api_pb.European `protobuf:"bytes,1,rep,name=europeans,proto3" json:"europeans,omitempty"`
}

func (s *europeanServiceServerImpl) ListEuropeans(ctx context.Context, req *api_pb.ListEuropeansRequest) (*api_pb.ListEuropeansResponse, error) {
	//	query := `query ListEuropeans($filter: string){
	query := `{
	Europeans(func: eq(type, "european")) %s {
		uid
		timestamp
		ticker
		strike
		und
		expiry
		putcall
	  }
	}`

	filter := ""
	if req.Ticker != "" && req.Timestamp != 0 {
		filter = fmt.Sprintf("@filter(eq(undticker, %s) AND eq(timestamp, %f))", req.Ticker, req.Timestamp)
	} else if req.Ticker != "" {
		filter = fmt.Sprintf("@filter(eq(undticker, %s))", req.Ticker)
	} else if req.Timestamp != 0 {
		filter = fmt.Sprintf("@filter(eq(timestamp, %f))", req.Timestamp)
	}

	query = fmt.Sprintf(query, filter)
	resp, err := s.db.NewTxn().Query(ctx, query)
	if err != nil {
		return nil, err
	}

	match := &api_pb.ListEuropeansResponse{}
	err = json.Unmarshal(resp.GetJson(), match)
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
	err = json.Unmarshal(resp.GetJson(), match)
	if err != nil {
		return nil, err
	}

	if len(match.Europeans) == 0 {
		return nil, nil
	}

	return match.Europeans[0], nil
}

// TODO: create european and add type in same tx
// TODO: create many europeans at once
func (s *europeanServiceServerImpl) CreateEuropean(ctx context.Context, req *api_pb.CreateEuropeanRequest) (*api_pb.European, error) {
	assigned, err := insertObj(ctx, s.db, req.European)
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	// TODO: get first UID rather than iterating through the map
	// maybe this goes away once we insert many contracts at once
	for k := range assigned.Uids {
		withType := typeDecorator{assigned.Uids[k], "european"}

		_, err = insertObj(ctx, s.db, withType)
		if err != nil {
			logrus.Errorln(err)
			return nil, err
		}

		req.European.Uid = withType.UID
	}

	if req.European.Uid == "" {
		err = fmt.Errorf("dgraph mutation returned no uid")
		logrus.Errorln(err)
		return nil, err
	}

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
