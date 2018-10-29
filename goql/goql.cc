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

#include "quantlib.h"

#include "spdlog/spdlog.h"
#include "spdlog/sinks/stdout_color_sinks.h"

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

const double putLBound = 0.2;
std::shared_ptr<spdlog::logger> console = spdlog::stdout_color_mt("goql");

// Logic and data behind the server's behavior.
class EuropeanOptionPricerServerImpl final : public EuropeanOptionPricer::Service {
        Status Price(ServerContext* context, const PriceRequest* request, PriceResponse* response) override {
                console->info("Incoming PriceRequest");
                Settings::instance().evaluationDate() = to_ql_date(request->pricingdate());
                boost::shared_ptr<BlackScholesMertonProcess> bsmProcess = buildBlackScholesMertonProcess(request->pricingdate(), request->marketdata());
                EuropeanOption europeanOption = buildEuropeanOption(request->contract().expiry(), request->contract().strike(), request->contract().putcall());
                europeanOption.setPricingEngine(boost::shared_ptr<PricingEngine>(new AnalyticEuropeanEngine(bsmProcess)));
                response->set_price(europeanOption.NPV());
                console->info("Outgoing PriceResponse");
                return Status::OK;
        }

        Status Greek(ServerContext* context, const GreekRequest* request,GreekResponse* response) override {
                console->info("Incoming GreekRequest");
                console->info("Outgoing GreekResponse");
                return Status::OK;
        }

        Status ImpliedVol(ServerContext* context, const ImpliedVolRequest* request, ImpliedVolResponse* response) override {
                console->info("Incoming ImpliedVolRequest");
                Settings::instance().evaluationDate() = to_ql_date(request->pricingdate());
                boost::shared_ptr<BlackScholesMertonProcess> bsmProcess = buildBlackScholesMertonProcess(request->pricingdate(), request->marketdata());
                auto s = request->marketdata().spot().index().value();
                for(int i = 0;i < request->quotes_size();i++) {
                        auto slice = request->quotes(i);

                        for(int k=0;k<slice.puts_size();k++){
                                auto put = slice.puts(k);
                                if (put.strike()/s > putLBound && put.strike()/s <= 1.0) {
                                        EuropeanOption europeanOption = buildEuropeanOption(slice.expiry(), put.strike(), "put");
                                        europeanOption.impliedVolatility(put.ask(), bsmProcess);
                                }
                        }

                        for(int k=0;k<slice.calls_size();k++){
                                auto call = slice.calls(k);
                                if (call.strike()/s > 1.0) {
                                        EuropeanOption europeanOption = buildEuropeanOption(slice.expiry(), call.strike(), "call");
                                        europeanOption.impliedVolatility(call.ask(), bsmProcess);
                                }
                        }
                }

                console->info("Outgoing ImpliedVolResponse");
                return Status::OK;
        }
};

std::string withServerAddress(int argc, char** argv)  {
        po::options_description desc("goql options");
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
        console->info("EuropeanOptionPricer grpc server ready on port at " + server_address);
        // Wait for the server to shutdown. Note that some other thread must be
        // responsible for shutting down the server for this call to ever return.
        server->Wait();

        return 0;
}
