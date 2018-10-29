#include <ql/quantlib.hpp>
#include <ql/time/calendars/target.hpp>
#include <ql/instruments/vanillaoption.hpp>
#include "boost/date_time/posix_time/posix_time.hpp"

#include "contract.pb.h"
#include "marketdata.pb.h"

using namespace pb;
using namespace std;
using namespace QuantLib;

#define to_ql_date(d) Date(boost::posix_time::from_time_t(d))

const double putLBound = 0.2;

boost::shared_ptr<BlackScholesMertonProcess> buildBlackScholesMertonProcess(int t0, const OptionMarket &mkt) {
    auto s = mkt.spot().index().value();
    auto v = mkt.vol().index().value();
    auto r = mkt.rate().index().value();

    Calendar calendar = TARGET();
    DayCounter dayCounter = Actual365Fixed();
    Date settlement = to_ql_date(t0)+0; //0 day settlement

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

    return boost::shared_ptr<BlackScholesMertonProcess>(
            new BlackScholesMertonProcess(
            underlyingH,
            flatDividendTS,
            flatTermStructure,
            flatVolTS));
}

EuropeanOption buildEuropeanOption(int expiry, double strike, string putcall) {
    auto optionType = static_cast<Option::Type>(putcall == "call" ? 1 : -1);
    boost::shared_ptr<StrikedTypePayoff> payoff(new PlainVanillaPayoff(optionType, strike));
    return EuropeanOption(payoff, boost::shared_ptr<Exercise>(new EuropeanExercise(to_ql_date(expiry))));
}

void calibrateImpliedVolSlice(const OptionQuoteSlice& quotes, const OptionMarket& market, boost::shared_ptr<BlackScholesMertonProcess> process, ImpliedVolSlice* calibratedSlice) {
        auto s = market.spot().index().value();
        calibratedSlice->set_expiry(quotes.expiry());
        calibratedSlice->set_timestamp(market.timestamp());

        for(int k=0;k<quotes.puts_size();k++){
                auto put = quotes.puts(k);
                if (put.strike()/s > putLBound && put.strike()/s <= 1.0) {
                        ImpliedVolQuote* ivQuote = calibratedSlice->add_quotes();
                        ivQuote->set_nbiteration(-1);
                        ivQuote->set_timestamp(put.timestamp());
                        ivQuote->set_allocated_input(new OptionQuote(put));

                        try {
                                EuropeanOption europeanOption = buildEuropeanOption(quotes.expiry(), put.strike(), "put");
                                ivQuote->set_vol(europeanOption.impliedVolatility(put.ask(), process));
                        } catch(exception& e) {
                                ivQuote->set_error(e.what());
                                calibratedSlice->set_iserror(true);
                        } catch(...) {
                                ivQuote->set_error("Unknown error when calling QuantLib");
                                calibratedSlice->set_iserror(true);
                        }
                }
        }

        for(int k=0;k<quotes.calls_size();k++){
                auto call = quotes.calls(k);
                if (call.strike()/s > 1.0) {
                        ImpliedVolQuote* ivQuote = calibratedSlice->add_quotes();
                        ivQuote->set_nbiteration(-1);
                        ivQuote->set_timestamp(call.timestamp());
                        ivQuote->set_allocated_input(new OptionQuote(call));

                        try {
                                EuropeanOption europeanOption = buildEuropeanOption(quotes.expiry(), call.strike(), "call");
                                ivQuote->set_vol(europeanOption.impliedVolatility(call.ask(), process));
                        } catch(exception& e) {
                                ivQuote->set_error(e.what());
                                calibratedSlice->set_iserror(true);
                        } catch(...) {
                                ivQuote->set_error("Unknown error when calling QuantLib");
                                calibratedSlice->set_iserror(true);
                        }
                }
        }
}
