package query

var (

	// ImpliedvolRequest represents a dgraph query returning data for a price request
	ImpliedvolRequest = `
{
	marketdata(func: eq(timestamp, $timestamp}})) @cascade { 
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

	quotes(func: eq(timestamp, $timestamp)) @cascade { 
		expiry
		puts (orderasc: strike){
			expand(_all_)  
		}
		calls (orderasc: strike) {
			expand(_all_)  
		}
	} 
}

fragment indexInfo {
	index @filter(eq(ticker, "$undTicker") or eq(ticker, "$rateTicker")) {
		timestamp
		ticker
		value
	}
}`
)
