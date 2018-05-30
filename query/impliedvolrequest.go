package query

var (

	// ImpliedvolRequest represents a dgraph query returning data for a price request
	ImpliedvolRequest = `
query ImpliedVolRequest($timestamp: float, $undTicker: string, $rateTicker: string){
	marketdata(func: eq(timestamp, $timestamp)) @cascade {
		spot {
			index @filter(eq(ticker, $undTicker)) {
				expand(_all_)  
			}
		}
		vol  {
			index @filter(eq(ticker, $undTicker)) {
				expand(_all_)  
			}
		}
		rate  {
			index @filter(eq(ticker, $rateTicker)) {
				expand(_all_)  
			}
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
}`
)
