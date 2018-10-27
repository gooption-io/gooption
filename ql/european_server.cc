#include <memory>
#include <string>
#include <iostream>

#include <grpc++/grpc++.h>
#include "service.pb.h"

#include <ql/quantlib.hpp>
#include <ql/time/calendars/target.hpp>
#include <ql/utilities/dataparsers.hpp>
#include <ql/instruments/vanillaoption.hpp>
#include "boost/date_time/posix_time/posix_time.hpp"
#include <ql/pricingengines/vanilla/analyticeuropeanengine.hpp>

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using proto::PriceRequest;
using proto::PriceResponse;
using proto::GreekRequest;
using proto::GreekResponse;
using proto::ImpliedVolRequest;
using proto::ImpliedVolResponse;
using proto::Gobs;

#ifdef BOOST_MSVC
/* Uncomment the following lines to unmask floating-point
   exceptions. Warning: unpredictable results can arise...
   See http://www.wilmott.com/messageview.cfm?catid=10&threadid=9481
   Is there anyone with a definitive word about this?
*/
// #include <float.h>
// namespace { unsigned int u = _controlfp(_EM_INEXACT, _MCW_EM); }
#endif

#include <iomanip>

using namespace std;
using namespace QuantLib;

#if defined(QL_ENABLE_SESSIONS)
namespace QuantLib {

    Integer sessionId() { return 0; }

}
#endif

#define to_ql_date(d) Date(boost::posix_time::from_time_t(d))

// Logic and data behind the server's behavior.
class GobsServerImpl final : public Gobs::Service {
        Status Price(ServerContext* context, const PriceRequest* request, PriceResponse* response) override {
                auto t0 = request->pricingdate();
                auto t = request->contract()->expiry();
                auto k = request->contract()->strike();
                auto s = request->marketdata()->spot()->index()->value();
                auto v = request->marketdata()->vol()->index()->value();
                auto r = request->marketdata()->rate()->index()->value();

                Calendar calendar = TARGET();
                DayCounter dayCounter = Actual365Fixed();
                Date settlement = to_ql_date(t0)+0; //0 day settlement
                Settings::instance().evaluationDate() = to_ql_date(t0);

                // bootstrap the yield/dividend/vol curves
                Handle<Quote> underlyingH(
                        boost::shared_ptr<Quote>(new SimpleQuote(s)));

                Handle<YieldTermStructure> flatTermStructure(
                        boost::shared_ptr<YieldTermStructure>(
                        new FlatForward(settlement, r, dayCounter)));

                Handle<YieldTermStructure> flatDividendTS(
                        boost::shared_ptr<YieldTermStructure>(
                        new FlatForward(settlement, 0., dayCounter)));

                Handle<BlackVolTermStructure> flatVolTS(
                        boost::shared_ptr<BlackVolTermStructure>(
                        new BlackConstantVol(settlement, calendar, v, dayCounter)));

                boost::shared_ptr<BlackScholesMertonProcess> bsmProcess(
                        new BlackScholesMertonProcess(
                        underlyingH,
                        flatDividendTS,
                        flatTermStructure,
                        flatVolTS));

                boost::shared_ptr<Exercise> europeanExercise(
                        new EuropeanExercise(to_ql_date(t)));

                boost::shared_ptr<StrikedTypePayoff> payoff(
                        new PlainVanillaPayoff(static_cast<Option::Type>(putcall),k));

                VanillaOption europeanOption(payoff, europeanExercise);
                europeanOption.setPricingEngine(
                        boost::shared_ptr<PricingEngine>(
                                new AnalyticEuropeanEngine(bsmProcess)));

                return Status::OK;
        }

        Status Greek(ServerContext* context, const GreekRequest* request,GreekResponse* response) override {
        return Status::OK;
        }

        Status ImpliedVol(ServerContext* context, const ImpliedVolRequest* request, ImpliedVolResponse* response) override {
        return Status::OK;
        }
};

void RunServer() {
  std::string server_address("0.0.0.0:50051");
  GobsServerImpl service;
  ServerBuilder builder;
  // Listen on the given address without any authentication mechanism.
  builder.AddListeningPort(server_address, grpc::InsecureServerCredentials());
  // Register "service" as the instance through which we'll communicate with
  // clients. In this case it corresponds to an *synchronous* service.
  builder.RegisterService(&service);
  // Finally assemble the server.
  std::unique_ptr<Server> server(builder.BuildAndStart());
  std::cout << "Server listening on " << server_address << std::endl;
  // Wait for the server to shutdown. Note that some other thread must be
  // responsible for shutting down the server for this call to ever return.
  server->Wait();
}

int main(int argc, char** argv) {
  RunServer();
  return 0;
}

// // Price a european option using FlatVol
// extern void EuropeanFlatVol(const PriceRequest* request, PriceResponse* response) {

//         auto t0 = request->pricingdate();
//         auto t = request->contract()->expiry();
//         auto k = request->contract()->strike();
//         auto s = request->marketdata()->spot()->index()->value();
//         auto v = request->marketdata()->vol()->index()->value();
//         auto r = request->marketdata()->rate()->index()->value();

//         Calendar calendar = TARGET();
//         DayCounter dayCounter = Actual365Fixed();
//         Date settlement = to_ql_date(t0)+0; //0 day settlement
//         Settings::instance().evaluationDate() = to_ql_date(t0);

//         // bootstrap the yield/dividend/vol curves
//         Handle<Quote> underlyingH(
//                 boost::shared_ptr<Quote>(new SimpleQuote(s)));

//         Handle<YieldTermStructure> flatTermStructure(
//                 boost::shared_ptr<YieldTermStructure>(
//                 new FlatForward(settlement, r, dayCounter)));

//         Handle<YieldTermStructure> flatDividendTS(
//                 boost::shared_ptr<YieldTermStructure>(
//                 new FlatForward(settlement, 0., dayCounter)));

//         Handle<BlackVolTermStructure> flatVolTS(
//                 boost::shared_ptr<BlackVolTermStructure>(
//                 new BlackConstantVol(settlement, calendar, v, dayCounter)));

//         boost::shared_ptr<BlackScholesMertonProcess> bsmProcess(
//                 new BlackScholesMertonProcess(
//                 underlyingH,
//                 flatDividendTS,
//                 flatTermStructure,
//                 flatVolTS));

//         boost::shared_ptr<Exercise> europeanExercise(
//                 new EuropeanExercise(to_ql_date(t)));

//         boost::shared_ptr<StrikedTypePayoff> payoff(
//                 new PlainVanillaPayoff(static_cast<Option::Type>(putcall),k));

//         VanillaOption europeanOption(payoff, europeanExercise);
//         europeanOption.setPricingEngine(
//                 boost::shared_ptr<PricingEngine>(
//                         new AnalyticEuropeanEngine(bsmProcess)));

//         return europeanOption.NPV();
// }