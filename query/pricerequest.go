package query

import (
	"github.com/valyala/fasttemplate"
)

var (

	// Query with optional tags goes here
	pricerequest = `
{
	contract(func: eq(ticker, "{{optionTicker}}")){
		ticker
		strike
		und as undticker
		expiry
		putcall
	}

	marketdata(func: eq(timestamp, {{timestamp}})) @cascade {
		spot {
			...indexInfo
		}
		vol  {
			...indexInfo
		}
		rate  {
			...indexInfo
		}
	}
}

fragment indexInfo {
	index @filter(eq(ticker, val(und)) or eq(ticker, "{{rateTicker}}")) {
		timestamp
		ticker
		value
	}
}`

	pricerequestTemplate = fasttemplate.New(pricerequest, startTag, endTag)
)

func init() {
	// For template reflection
	AllTemplates["pricerequest"] = pricerequestTemplate
}

func GetPriceRequestQuery(timestamp, optionTicker, rateTicker string) string {

	return pricerequestTemplate.ExecuteString(map[string]interface{}{
		"timestamp": timestamp, "optionTicker": optionTicker, "rateTicker": rateTicker,
	})
}
