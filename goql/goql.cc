#include <memory>
#include <string>
#include <iostream>
#include <iomanip>

#include <grpc++/grpc++.h>
#include "service.grpc.pb.h"
#include "marketdata.pb.h"

#include <ql/quantlib.hpp>
#include <ql/time/calendars/target.hpp>
#include <ql/utilities/dataparsers.hpp>
#include <ql/instruments/vanillaoption.hpp>
#include <ql/pricingengines/vanilla/analyticeuropeanengine.hpp>

#include "boost/program_options.hpp"
#include "boost/date_time/posix_time/posix_time.hpp"

#include "ql.cc"

#include "spdlog/spdlog.h"
#include "spdlog/sinks/stdout_color_sinks.h"

using namespace pb;
using namespace std;
using namespace grpc;
using namespace QuantLib;
namespace po = boost::program_options;

#if defined(QL_ENABLE_SESSIONS)
namespace QuantLib {
Integer sessionId() { return 0; }
}
#endif

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
        auto priceRequest = request->request();
        Settings::instance().evaluationDate() = to_ql_date(priceRequest.pricingdate());
        boost::shared_ptr<BlackScholesMertonProcess> bsmProcess = buildBlackScholesMertonProcess(priceRequest.pricingdate(), priceRequest.marketdata());
        EuropeanOption europeanOption = buildEuropeanOption(priceRequest.contract().expiry(), priceRequest.contract().strike(), priceRequest.contract().putcall());
        europeanOption.setPricingEngine(boost::shared_ptr<PricingEngine>(new AnalyticEuropeanEngine(bsmProcess)));

        vector<string> greeks;
        for(int i = 0; i < request->greek_size();i++) {
            if(request->greek(i) == "all") {
                greeks = vector<string>{"delta", "gamma", "vega", "rho", "theta"};
                break;
            }
            greeks.push_back(request->greek(i));
        }

        for(int i = 0; i < greeks.size();i++) {
            if(greeks[i] == "delta") {
                auto delta = response->add_greeks();
                delta->set_label("delta");
                delta->set_value(europeanOption.delta());
            }
            if(greeks[i] == "gamma") {
                auto gamma = response->add_greeks();
                gamma->set_label("gamma");
                gamma->set_value(europeanOption.gamma());
            }
            if(greeks[i] == "vega") {
                auto vega = response->add_greeks();
                vega->set_label("vega");
                vega->set_value(europeanOption.vega());
            }
            if(greeks[i] == "rho") {
                auto rho = response->add_greeks();
                rho->set_label("rho");
                rho->set_value(europeanOption.rho());
            }
            if(greeks[i] == "theta") {
                auto theta = response->add_greeks();
                theta->set_label("theta");
                theta->set_value(europeanOption.theta());
            }
        }

        console->info("Outgoing GreekResponse");
        return Status::OK;
    }

    Status ImpliedVol(ServerContext* context, const ImpliedVolRequest* request, ImpliedVolResponse* response) override {
        console->info("Incoming ImpliedVolRequest");
        Settings::instance().evaluationDate() = to_ql_date(request->pricingdate());
        boost::shared_ptr<BlackScholesMertonProcess> bsmProcess = buildBlackScholesMertonProcess(request->pricingdate(), request->marketdata());

        for(int i = 0;i < request->quotes_size();i++) {
            console->debug("Calibrating slice");
            auto quotes = request->quotes(i);
            auto calibratedSlice = response->mutable_volsurface()->add_slices();
            calibrateImpliedVolSlice(quotes, request->marketdata(), bsmProcess, calibratedSlice);
            console->debug(calibratedSlice->ShortDebugString());
        }

        console->info("Outgoing ImpliedVolResponse");
        console->debug(response->ShortDebugString());
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
