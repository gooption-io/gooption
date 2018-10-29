#include <ql/quantlib.hpp>
#include <ql/time/calendars/target.hpp>
#include <ql/instruments/vanillaoption.hpp>
#include "boost/date_time/posix_time/posix_time.hpp"

#include "contract.pb.h"
#include "marketdata.pb.h"

using namespace std;
using namespace QuantLib;
using pb::European;
using pb::OptionMarket;

#define to_ql_date(d) Date(boost::posix_time::from_time_t(d))

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
