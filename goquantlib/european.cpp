#include "european.hpp"
#include <iostream>
#include <ql/quantlib.hpp>
#include <ql/instruments/vanillaoption.hpp>
#include <ql/pricingengines/vanilla/analyticeuropeanengine.hpp>
#include <ql/pricingengines/vanilla/fdeuropeanengine.hpp>
#include <ql/time/calendars/target.hpp>
#include <ql/utilities/dataparsers.hpp>
#include "boost/date_time/posix_time/posix_time.hpp"

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

// Price a european option using FlatVol
double EuropeanFlatVol(
    double s, 
    double r, 
    double q, 
    double v, 
    double k, 
    int t0, 
    int t, 
    int putcall) {

    Calendar calendar = TARGET();
    DayCounter dayCounter = Actual365Fixed();   
    Date settlement = to_ql_date(t)+0; //0 day settlement0
   
    // bootstrap the yield/dividend/vol curves
    Handle<Quote> underlyingH(boost::shared_ptr<Quote>(new SimpleQuote(s)));
    Handle<YieldTermStructure> flatTermStructure(
            boost::shared_ptr<YieldTermStructure>(
            new FlatForward(settlement, r, dayCounter)));

    Handle<YieldTermStructure> flatDividendTS(
            boost::shared_ptr<YieldTermStructure>(
            new FlatForward(settlement, q, dayCounter)));

    Handle<BlackVolTermStructure> flatVolTS(
            boost::shared_ptr<BlackVolTermStructure>(
            new BlackConstantVol(settlement, calendar, v, dayCounter)));

    boost::shared_ptr<BlackScholesMertonProcess> bsmProcess(
            new BlackScholesMertonProcess(
                underlyingH, 
                flatDividendTS,
                flatTermStructure, 
                flatVolTS));

    boost::shared_ptr<Exercise> europeanExercise(new EuropeanExercise(to_ql_date(t)));                    
    boost::shared_ptr<StrikedTypePayoff> payoff(new PlainVanillaPayoff(static_cast<Option::Type>(putcall),k));
    VanillaOption europeanOption(payoff, europeanExercise);
    europeanOption.setPricingEngine(boost::shared_ptr<PricingEngine>(new AnalyticEuropeanEngine(bsmProcess)));

    Settings::instance().evaluationDate() = to_ql_date(t0);  
    return europeanOption.NPV();
}