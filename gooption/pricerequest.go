package main

var (
	// PriceRequest represents a dgraph query returning data for a price request
	PriceRequest = `
query PriceRequest($timestamp: float, $optionTicker: string, $rateTicker: string){
	contract(func: eq(ticker, $optionTicker)){
		ticker
		strike
		und as undticker
		expiry
		putcall
	}

	marketdata(func: eq(timestamp, $timestamp)) @cascade {
		spot {
			index @filter(eq(ticker, val(und))) {
				timestamp
				ticker
				value
			}
		}
		vol  {
			index @filter(eq(ticker, val(und))) {
				timestamp
				ticker
				value
			}
		}
		rate  {
			index @filter(eq(ticker, $rateTicker)) {
				timestamp
				ticker
				value
			}
		}
	}
}`
)
