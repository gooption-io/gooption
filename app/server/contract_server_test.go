package server

import (
	"context"
	"testing"

	api_pb "github.com/gooption-io/gooption/api"
)

func Test_ContractServiceServer_ListContracts(t *testing.T) {
	svr := NewContractServiceServer()

	ctx := context.Background()
	req := &api_pb.ListContractsRequest{}

	resp, err := svr.ListContracts(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_ContractServiceServer_GetContract(t *testing.T) {
	svr := NewContractServiceServer()

	ctx := context.Background()
	req := &api_pb.GetContractRequest{}

	resp, err := svr.GetContract(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_ContractServiceServer_CreateContract(t *testing.T) {
	svr := NewContractServiceServer()

	ctx := context.Background()
	req := &api_pb.CreateContractRequest{}

	resp, err := svr.CreateContract(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_ContractServiceServer_UpdateContract(t *testing.T) {
	svr := NewContractServiceServer()

	ctx := context.Background()
	req := &api_pb.UpdateContractRequest{}

	resp, err := svr.UpdateContract(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_ContractServiceServer_DeleteContract(t *testing.T) {
	svr := NewContractServiceServer()

	ctx := context.Background()
	req := &api_pb.DeleteContractRequest{}

	resp, err := svr.DeleteContract(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}
