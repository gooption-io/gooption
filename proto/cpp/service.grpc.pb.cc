// Generated by the gRPC C++ plugin.
// If you make any local change, they will be lost.
// source: service.proto

#include "service.pb.h"
#include "service.grpc.pb.h"

#include <grpcpp/impl/codegen/async_stream.h>
#include <grpcpp/impl/codegen/async_unary_call.h>
#include <grpcpp/impl/codegen/channel_interface.h>
#include <grpcpp/impl/codegen/client_unary_call.h>
#include <grpcpp/impl/codegen/method_handler_impl.h>
#include <grpcpp/impl/codegen/rpc_service_method.h>
#include <grpcpp/impl/codegen/service_type.h>
#include <grpcpp/impl/codegen/sync_stream.h>
namespace proto {

static const char* Gobs_method_names[] = {
  "/proto.Gobs/Price",
  "/proto.Gobs/Greek",
  "/proto.Gobs/ImpliedVol",
};

std::unique_ptr< Gobs::Stub> Gobs::NewStub(const std::shared_ptr< ::grpc::ChannelInterface>& channel, const ::grpc::StubOptions& options) {
  (void)options;
  std::unique_ptr< Gobs::Stub> stub(new Gobs::Stub(channel));
  return stub;
}

Gobs::Stub::Stub(const std::shared_ptr< ::grpc::ChannelInterface>& channel)
  : channel_(channel), rpcmethod_Price_(Gobs_method_names[0], ::grpc::internal::RpcMethod::NORMAL_RPC, channel)
  , rpcmethod_Greek_(Gobs_method_names[1], ::grpc::internal::RpcMethod::NORMAL_RPC, channel)
  , rpcmethod_ImpliedVol_(Gobs_method_names[2], ::grpc::internal::RpcMethod::NORMAL_RPC, channel)
  {}

::grpc::Status Gobs::Stub::Price(::grpc::ClientContext* context, const ::proto::PriceRequest& request, ::proto::PriceResponse* response) {
  return ::grpc::internal::BlockingUnaryCall(channel_.get(), rpcmethod_Price_, context, request, response);
}

::grpc::ClientAsyncResponseReader< ::proto::PriceResponse>* Gobs::Stub::AsyncPriceRaw(::grpc::ClientContext* context, const ::proto::PriceRequest& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::proto::PriceResponse>::Create(channel_.get(), cq, rpcmethod_Price_, context, request, true);
}

::grpc::ClientAsyncResponseReader< ::proto::PriceResponse>* Gobs::Stub::PrepareAsyncPriceRaw(::grpc::ClientContext* context, const ::proto::PriceRequest& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::proto::PriceResponse>::Create(channel_.get(), cq, rpcmethod_Price_, context, request, false);
}

::grpc::Status Gobs::Stub::Greek(::grpc::ClientContext* context, const ::proto::GreekRequest& request, ::proto::GreekResponse* response) {
  return ::grpc::internal::BlockingUnaryCall(channel_.get(), rpcmethod_Greek_, context, request, response);
}

::grpc::ClientAsyncResponseReader< ::proto::GreekResponse>* Gobs::Stub::AsyncGreekRaw(::grpc::ClientContext* context, const ::proto::GreekRequest& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::proto::GreekResponse>::Create(channel_.get(), cq, rpcmethod_Greek_, context, request, true);
}

::grpc::ClientAsyncResponseReader< ::proto::GreekResponse>* Gobs::Stub::PrepareAsyncGreekRaw(::grpc::ClientContext* context, const ::proto::GreekRequest& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::proto::GreekResponse>::Create(channel_.get(), cq, rpcmethod_Greek_, context, request, false);
}

::grpc::Status Gobs::Stub::ImpliedVol(::grpc::ClientContext* context, const ::proto::ImpliedVolRequest& request, ::proto::ImpliedVolResponse* response) {
  return ::grpc::internal::BlockingUnaryCall(channel_.get(), rpcmethod_ImpliedVol_, context, request, response);
}

::grpc::ClientAsyncResponseReader< ::proto::ImpliedVolResponse>* Gobs::Stub::AsyncImpliedVolRaw(::grpc::ClientContext* context, const ::proto::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::proto::ImpliedVolResponse>::Create(channel_.get(), cq, rpcmethod_ImpliedVol_, context, request, true);
}

::grpc::ClientAsyncResponseReader< ::proto::ImpliedVolResponse>* Gobs::Stub::PrepareAsyncImpliedVolRaw(::grpc::ClientContext* context, const ::proto::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) {
  return ::grpc::internal::ClientAsyncResponseReaderFactory< ::proto::ImpliedVolResponse>::Create(channel_.get(), cq, rpcmethod_ImpliedVol_, context, request, false);
}

Gobs::Service::Service() {
  AddMethod(new ::grpc::internal::RpcServiceMethod(
      Gobs_method_names[0],
      ::grpc::internal::RpcMethod::NORMAL_RPC,
      new ::grpc::internal::RpcMethodHandler< Gobs::Service, ::proto::PriceRequest, ::proto::PriceResponse>(
          std::mem_fn(&Gobs::Service::Price), this)));
  AddMethod(new ::grpc::internal::RpcServiceMethod(
      Gobs_method_names[1],
      ::grpc::internal::RpcMethod::NORMAL_RPC,
      new ::grpc::internal::RpcMethodHandler< Gobs::Service, ::proto::GreekRequest, ::proto::GreekResponse>(
          std::mem_fn(&Gobs::Service::Greek), this)));
  AddMethod(new ::grpc::internal::RpcServiceMethod(
      Gobs_method_names[2],
      ::grpc::internal::RpcMethod::NORMAL_RPC,
      new ::grpc::internal::RpcMethodHandler< Gobs::Service, ::proto::ImpliedVolRequest, ::proto::ImpliedVolResponse>(
          std::mem_fn(&Gobs::Service::ImpliedVol), this)));
}

Gobs::Service::~Service() {
}

::grpc::Status Gobs::Service::Price(::grpc::ServerContext* context, const ::proto::PriceRequest* request, ::proto::PriceResponse* response) {
  (void) context;
  (void) request;
  (void) response;
  return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
}

::grpc::Status Gobs::Service::Greek(::grpc::ServerContext* context, const ::proto::GreekRequest* request, ::proto::GreekResponse* response) {
  (void) context;
  (void) request;
  (void) response;
  return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
}

::grpc::Status Gobs::Service::ImpliedVol(::grpc::ServerContext* context, const ::proto::ImpliedVolRequest* request, ::proto::ImpliedVolResponse* response) {
  (void) context;
  (void) request;
  (void) response;
  return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
}


}  // namespace proto

