// Code generated by github.com/izumin5210/grapi. DO NOT EDIT.

package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	api_pb "github.com/gooption-io/gooption/api"
)

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *stockServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	api_pb.RegisterStockServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *stockServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return api_pb.RegisterStockServiceHandler(ctx, mux, conn)
}
