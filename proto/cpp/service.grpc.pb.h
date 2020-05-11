// Generated by the gRPC C++ plugin.
// If you make any local change, they will be lost.
// source: service.proto
#ifndef GRPC_service_2eproto__INCLUDED
#define GRPC_service_2eproto__INCLUDED

#include "service.pb.h"

#include <grpcpp/impl/codegen/async_generic_service.h>
#include <grpcpp/impl/codegen/async_stream.h>
#include <grpcpp/impl/codegen/async_unary_call.h>
#include <grpcpp/impl/codegen/method_handler_impl.h>
#include <grpcpp/impl/codegen/proto_utils.h>
#include <grpcpp/impl/codegen/rpc_method.h>
#include <grpcpp/impl/codegen/service_type.h>
#include <grpcpp/impl/codegen/status.h>
#include <grpcpp/impl/codegen/stub_options.h>
#include <grpcpp/impl/codegen/sync_stream.h>

namespace grpc {
class CompletionQueue;
class Channel;
class ServerCompletionQueue;
class ServerContext;
}  // namespace grpc

namespace pb {

class EuropeanOptionPricer final {
 public:
  static constexpr char const* service_full_name() {
    return "pb.EuropeanOptionPricer";
  }
  class StubInterface {
   public:
    virtual ~StubInterface() {}
    virtual ::grpc::Status Price(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::pb::PriceResponse* response) = 0;
    std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::PriceResponse>> AsyncPrice(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::PriceResponse>>(AsyncPriceRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::PriceResponse>> PrepareAsyncPrice(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::PriceResponse>>(PrepareAsyncPriceRaw(context, request, cq));
    }
    virtual ::grpc::Status Greek(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::pb::GreekResponse* response) = 0;
    std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::GreekResponse>> AsyncGreek(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::GreekResponse>>(AsyncGreekRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::GreekResponse>> PrepareAsyncGreek(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::GreekResponse>>(PrepareAsyncGreekRaw(context, request, cq));
    }
    virtual ::grpc::Status ImpliedVol(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::pb::ImpliedVolResponse* response) = 0;
    std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::ImpliedVolResponse>> AsyncImpliedVol(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::ImpliedVolResponse>>(AsyncImpliedVolRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::ImpliedVolResponse>> PrepareAsyncImpliedVol(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReaderInterface< ::pb::ImpliedVolResponse>>(PrepareAsyncImpliedVolRaw(context, request, cq));
    }
  private:
    virtual ::grpc::ClientAsyncResponseReaderInterface< ::pb::PriceResponse>* AsyncPriceRaw(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::grpc::CompletionQueue* cq) = 0;
    virtual ::grpc::ClientAsyncResponseReaderInterface< ::pb::PriceResponse>* PrepareAsyncPriceRaw(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::grpc::CompletionQueue* cq) = 0;
    virtual ::grpc::ClientAsyncResponseReaderInterface< ::pb::GreekResponse>* AsyncGreekRaw(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::grpc::CompletionQueue* cq) = 0;
    virtual ::grpc::ClientAsyncResponseReaderInterface< ::pb::GreekResponse>* PrepareAsyncGreekRaw(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::grpc::CompletionQueue* cq) = 0;
    virtual ::grpc::ClientAsyncResponseReaderInterface< ::pb::ImpliedVolResponse>* AsyncImpliedVolRaw(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) = 0;
    virtual ::grpc::ClientAsyncResponseReaderInterface< ::pb::ImpliedVolResponse>* PrepareAsyncImpliedVolRaw(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) = 0;
  };
  class Stub final : public StubInterface {
   public:
    Stub(const std::shared_ptr< ::grpc::ChannelInterface>& channel);
    ::grpc::Status Price(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::pb::PriceResponse* response) override;
    std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::PriceResponse>> AsyncPrice(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::PriceResponse>>(AsyncPriceRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::PriceResponse>> PrepareAsyncPrice(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::PriceResponse>>(PrepareAsyncPriceRaw(context, request, cq));
    }
    ::grpc::Status Greek(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::pb::GreekResponse* response) override;
    std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::GreekResponse>> AsyncGreek(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::GreekResponse>>(AsyncGreekRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::GreekResponse>> PrepareAsyncGreek(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::GreekResponse>>(PrepareAsyncGreekRaw(context, request, cq));
    }
    ::grpc::Status ImpliedVol(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::pb::ImpliedVolResponse* response) override;
    std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::ImpliedVolResponse>> AsyncImpliedVol(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::ImpliedVolResponse>>(AsyncImpliedVolRaw(context, request, cq));
    }
    std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::ImpliedVolResponse>> PrepareAsyncImpliedVol(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) {
      return std::unique_ptr< ::grpc::ClientAsyncResponseReader< ::pb::ImpliedVolResponse>>(PrepareAsyncImpliedVolRaw(context, request, cq));
    }

   private:
    std::shared_ptr< ::grpc::ChannelInterface> channel_;
    ::grpc::ClientAsyncResponseReader< ::pb::PriceResponse>* AsyncPriceRaw(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::grpc::CompletionQueue* cq) override;
    ::grpc::ClientAsyncResponseReader< ::pb::PriceResponse>* PrepareAsyncPriceRaw(::grpc::ClientContext* context, const ::pb::PriceRequest& request, ::grpc::CompletionQueue* cq) override;
    ::grpc::ClientAsyncResponseReader< ::pb::GreekResponse>* AsyncGreekRaw(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::grpc::CompletionQueue* cq) override;
    ::grpc::ClientAsyncResponseReader< ::pb::GreekResponse>* PrepareAsyncGreekRaw(::grpc::ClientContext* context, const ::pb::GreekRequest& request, ::grpc::CompletionQueue* cq) override;
    ::grpc::ClientAsyncResponseReader< ::pb::ImpliedVolResponse>* AsyncImpliedVolRaw(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) override;
    ::grpc::ClientAsyncResponseReader< ::pb::ImpliedVolResponse>* PrepareAsyncImpliedVolRaw(::grpc::ClientContext* context, const ::pb::ImpliedVolRequest& request, ::grpc::CompletionQueue* cq) override;
    const ::grpc::internal::RpcMethod rpcmethod_Price_;
    const ::grpc::internal::RpcMethod rpcmethod_Greek_;
    const ::grpc::internal::RpcMethod rpcmethod_ImpliedVol_;
  };
  static std::unique_ptr<Stub> NewStub(const std::shared_ptr< ::grpc::ChannelInterface>& channel, const ::grpc::StubOptions& options = ::grpc::StubOptions());

  class Service : public ::grpc::Service {
   public:
    Service();
    virtual ~Service();
    virtual ::grpc::Status Price(::grpc::ServerContext* context, const ::pb::PriceRequest* request, ::pb::PriceResponse* response);
    virtual ::grpc::Status Greek(::grpc::ServerContext* context, const ::pb::GreekRequest* request, ::pb::GreekResponse* response);
    virtual ::grpc::Status ImpliedVol(::grpc::ServerContext* context, const ::pb::ImpliedVolRequest* request, ::pb::ImpliedVolResponse* response);
  };
  template <class BaseClass>
  class WithAsyncMethod_Price : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithAsyncMethod_Price() {
      ::grpc::Service::MarkMethodAsync(0);
    }
    ~WithAsyncMethod_Price() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status Price(::grpc::ServerContext* context, const ::pb::PriceRequest* request, ::pb::PriceResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestPrice(::grpc::ServerContext* context, ::pb::PriceRequest* request, ::grpc::ServerAsyncResponseWriter< ::pb::PriceResponse>* response, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncUnary(0, context, request, response, new_call_cq, notification_cq, tag);
    }
  };
  template <class BaseClass>
  class WithAsyncMethod_Greek : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithAsyncMethod_Greek() {
      ::grpc::Service::MarkMethodAsync(1);
    }
    ~WithAsyncMethod_Greek() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status Greek(::grpc::ServerContext* context, const ::pb::GreekRequest* request, ::pb::GreekResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestGreek(::grpc::ServerContext* context, ::pb::GreekRequest* request, ::grpc::ServerAsyncResponseWriter< ::pb::GreekResponse>* response, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncUnary(1, context, request, response, new_call_cq, notification_cq, tag);
    }
  };
  template <class BaseClass>
  class WithAsyncMethod_ImpliedVol : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithAsyncMethod_ImpliedVol() {
      ::grpc::Service::MarkMethodAsync(2);
    }
    ~WithAsyncMethod_ImpliedVol() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status ImpliedVol(::grpc::ServerContext* context, const ::pb::ImpliedVolRequest* request, ::pb::ImpliedVolResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestImpliedVol(::grpc::ServerContext* context, ::pb::ImpliedVolRequest* request, ::grpc::ServerAsyncResponseWriter< ::pb::ImpliedVolResponse>* response, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncUnary(2, context, request, response, new_call_cq, notification_cq, tag);
    }
  };
  typedef WithAsyncMethod_Price<WithAsyncMethod_Greek<WithAsyncMethod_ImpliedVol<Service > > > AsyncService;
  template <class BaseClass>
  class WithGenericMethod_Price : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithGenericMethod_Price() {
      ::grpc::Service::MarkMethodGeneric(0);
    }
    ~WithGenericMethod_Price() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status Price(::grpc::ServerContext* context, const ::pb::PriceRequest* request, ::pb::PriceResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
  };
  template <class BaseClass>
  class WithGenericMethod_Greek : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithGenericMethod_Greek() {
      ::grpc::Service::MarkMethodGeneric(1);
    }
    ~WithGenericMethod_Greek() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status Greek(::grpc::ServerContext* context, const ::pb::GreekRequest* request, ::pb::GreekResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
  };
  template <class BaseClass>
  class WithGenericMethod_ImpliedVol : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithGenericMethod_ImpliedVol() {
      ::grpc::Service::MarkMethodGeneric(2);
    }
    ~WithGenericMethod_ImpliedVol() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status ImpliedVol(::grpc::ServerContext* context, const ::pb::ImpliedVolRequest* request, ::pb::ImpliedVolResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
  };
  template <class BaseClass>
  class WithRawMethod_Price : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithRawMethod_Price() {
      ::grpc::Service::MarkMethodRaw(0);
    }
    ~WithRawMethod_Price() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status Price(::grpc::ServerContext* context, const ::pb::PriceRequest* request, ::pb::PriceResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestPrice(::grpc::ServerContext* context, ::grpc::ByteBuffer* request, ::grpc::ServerAsyncResponseWriter< ::grpc::ByteBuffer>* response, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncUnary(0, context, request, response, new_call_cq, notification_cq, tag);
    }
  };
  template <class BaseClass>
  class WithRawMethod_Greek : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithRawMethod_Greek() {
      ::grpc::Service::MarkMethodRaw(1);
    }
    ~WithRawMethod_Greek() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status Greek(::grpc::ServerContext* context, const ::pb::GreekRequest* request, ::pb::GreekResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestGreek(::grpc::ServerContext* context, ::grpc::ByteBuffer* request, ::grpc::ServerAsyncResponseWriter< ::grpc::ByteBuffer>* response, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncUnary(1, context, request, response, new_call_cq, notification_cq, tag);
    }
  };
  template <class BaseClass>
  class WithRawMethod_ImpliedVol : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithRawMethod_ImpliedVol() {
      ::grpc::Service::MarkMethodRaw(2);
    }
    ~WithRawMethod_ImpliedVol() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable synchronous version of this method
    ::grpc::Status ImpliedVol(::grpc::ServerContext* context, const ::pb::ImpliedVolRequest* request, ::pb::ImpliedVolResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    void RequestImpliedVol(::grpc::ServerContext* context, ::grpc::ByteBuffer* request, ::grpc::ServerAsyncResponseWriter< ::grpc::ByteBuffer>* response, ::grpc::CompletionQueue* new_call_cq, ::grpc::ServerCompletionQueue* notification_cq, void *tag) {
      ::grpc::Service::RequestAsyncUnary(2, context, request, response, new_call_cq, notification_cq, tag);
    }
  };
  template <class BaseClass>
  class WithStreamedUnaryMethod_Price : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithStreamedUnaryMethod_Price() {
      ::grpc::Service::MarkMethodStreamed(0,
        new ::grpc::internal::StreamedUnaryHandler< ::pb::PriceRequest, ::pb::PriceResponse>(std::bind(&WithStreamedUnaryMethod_Price<BaseClass>::StreamedPrice, this, std::placeholders::_1, std::placeholders::_2)));
    }
    ~WithStreamedUnaryMethod_Price() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable regular version of this method
    ::grpc::Status Price(::grpc::ServerContext* context, const ::pb::PriceRequest* request, ::pb::PriceResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    // replace default version of method with streamed unary
    virtual ::grpc::Status StreamedPrice(::grpc::ServerContext* context, ::grpc::ServerUnaryStreamer< ::pb::PriceRequest,::pb::PriceResponse>* server_unary_streamer) = 0;
  };
  template <class BaseClass>
  class WithStreamedUnaryMethod_Greek : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithStreamedUnaryMethod_Greek() {
      ::grpc::Service::MarkMethodStreamed(1,
        new ::grpc::internal::StreamedUnaryHandler< ::pb::GreekRequest, ::pb::GreekResponse>(std::bind(&WithStreamedUnaryMethod_Greek<BaseClass>::StreamedGreek, this, std::placeholders::_1, std::placeholders::_2)));
    }
    ~WithStreamedUnaryMethod_Greek() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable regular version of this method
    ::grpc::Status Greek(::grpc::ServerContext* context, const ::pb::GreekRequest* request, ::pb::GreekResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    // replace default version of method with streamed unary
    virtual ::grpc::Status StreamedGreek(::grpc::ServerContext* context, ::grpc::ServerUnaryStreamer< ::pb::GreekRequest,::pb::GreekResponse>* server_unary_streamer) = 0;
  };
  template <class BaseClass>
  class WithStreamedUnaryMethod_ImpliedVol : public BaseClass {
   private:
    void BaseClassMustBeDerivedFromService(const Service *service) {}
   public:
    WithStreamedUnaryMethod_ImpliedVol() {
      ::grpc::Service::MarkMethodStreamed(2,
        new ::grpc::internal::StreamedUnaryHandler< ::pb::ImpliedVolRequest, ::pb::ImpliedVolResponse>(std::bind(&WithStreamedUnaryMethod_ImpliedVol<BaseClass>::StreamedImpliedVol, this, std::placeholders::_1, std::placeholders::_2)));
    }
    ~WithStreamedUnaryMethod_ImpliedVol() override {
      BaseClassMustBeDerivedFromService(this);
    }
    // disable regular version of this method
    ::grpc::Status ImpliedVol(::grpc::ServerContext* context, const ::pb::ImpliedVolRequest* request, ::pb::ImpliedVolResponse* response) override {
      abort();
      return ::grpc::Status(::grpc::StatusCode::UNIMPLEMENTED, "");
    }
    // replace default version of method with streamed unary
    virtual ::grpc::Status StreamedImpliedVol(::grpc::ServerContext* context, ::grpc::ServerUnaryStreamer< ::pb::ImpliedVolRequest,::pb::ImpliedVolResponse>* server_unary_streamer) = 0;
  };
  typedef WithStreamedUnaryMethod_Price<WithStreamedUnaryMethod_Greek<WithStreamedUnaryMethod_ImpliedVol<Service > > > StreamedUnaryService;
  typedef Service SplitStreamedService;
  typedef WithStreamedUnaryMethod_Price<WithStreamedUnaryMethod_Greek<WithStreamedUnaryMethod_ImpliedVol<Service > > > StreamedService;
};

}  // namespace pb


#endif  // GRPC_service_2eproto__INCLUDED