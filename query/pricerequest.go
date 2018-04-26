package query

var (

	// PriceRequest represents a dgraph query returning data for a price request
	PriceRequest = `
{
	contract(func: eq(ticker, "$optionTicker")){
		ticker
		strike
		und as undticker
		expiry
		putcall
	}

	marketdata(func: eq(timestamp, $timestamp)) @cascade {
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
	index @filter(eq(ticker, val(und)) or eq(ticker, "$rateTicker")) {
		timestamp
		ticker
		value
	}
}`
)
