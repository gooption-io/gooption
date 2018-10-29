#include <memory>
#include <string>
#include <iostream>
#include <iomanip>

#include <grpc++/grpc++.h>
#include "service.grpc.pb.h"

#include <ql/quantlib.hpp>
#include <ql/time/calendars/target.hpp>
#include <ql/utilities/dataparsers.hpp>
#include <ql/instruments/vanillaoption.hpp>
#include <ql/pricingengines/vanilla/analyticeuropeanengine.hpp>

#include "boost/program_options.hpp"
#include "boost/date_time/posix_time/posix_time.hpp"

using namespace std;
using namespace QuantLib;
namespace po = boost::program_options;

using grpc::Server;
using grpc::ServerBuilder;
using grpc::ServerContext;
using grpc::Status;
using pb::PriceRequest;
using pb::PriceResponse;
using pb::GreekRequest;
using pb::GreekResponse;
using pb::ImpliedVolRequest;
using pb::ImpliedVolResponse;
using pb::EuropeanOptionPricer;

#if defined(QL_ENABLE_SESSIONS)
namespace QuantLib {
    Integer sessionId() { return 0; }
}
#endif

#define to_ql_date(d) Date(boost::posix_time::from_time_t(d))

// Logic and data behind the server's behavior.
class EuropeanOptionPricerServerImpl final : public EuropeanOptionPricer::Service {
        Status Price(ServerContext* context, const PriceRequest* request, PriceResponse* response) override {
                int t0 = request->pricingdate();
                int t = request->contract().expiry();

                auto k = request->contract().strike();
                auto putcall = static_cast<Option::Type>(request->contract().putcall() == "call" ? 1 : -1);

                auto s = request->marketdata().spot().index().value();
                auto v = request->marketdata().vol().index().value();
                auto r = request->marketdata().rate().index().value();

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
                        new PlainVanillaPayoff(putcall,k));

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

std::string withServerAddress(int argc, char** argv)  {
        po::options_description desc("Allowed options");
        desc.add_options()("tcp-listen-address", po::value<string>(), "tcp port");

        po::variables_map vm;
        po::store(po::parse_command_line(argc, argv, desc), vm);
        po::notify(vm);

        if (!(vm.count("tcp-listen-address"))) {
                return ":50051";
        }

        return vm["tcp-listen-address"].as<std::string>();
}

int main(int argc, char** argv) {
        ServerBuilder builder;
        EuropeanOptionPricerServerImpl service;

        // Listen on the given address without any authentication mechanism.
        std::string server_address = withServerAddress(argc, argv);
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

        return 0;
}
