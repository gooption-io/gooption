package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "github.com/gooption-io/gooption/api"
)

// ContractServiceServer is a composite interface of api_pb.ContractServiceServer and grapiserver.Server.
type ContractServiceServer interface {
	api_pb.ContractServiceServer
	grapiserver.Server
}

// NewContractServiceServer creates a new ContractServiceServer instance.
func NewContractServiceServer() ContractServiceServer {
	return &contractServiceServerImpl{}
}

type contractServiceServerImpl struct {
}

func (s *contractServiceServerImpl) ListContracts(ctx context.Context, req *api_pb.ListContractsRequest) (*api_pb.ListContractsResponse, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *contractServiceServerImpl) GetContract(ctx context.Context, req *api_pb.GetContractRequest) (*api_pb.Contract, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *contractServiceServerImpl) CreateContract(ctx context.Context, req *api_pb.CreateContractRequest) (*api_pb.Contract, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *contractServiceServerImpl) UpdateContract(ctx context.Context, req *api_pb.UpdateContractRequest) (*api_pb.Contract, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *contractServiceServerImpl) DeleteContract(ctx context.Context, req *api_pb.DeleteContractRequest) (*empty.Empty, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
