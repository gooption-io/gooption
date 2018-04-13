package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetImpliedVolRequestQuery(t *testing.T) {
	expected := `
{
	marketdata(func: eq(timestamp, 1513551151)) @cascade { 
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

	quotes(func: eq(timestamp, 1513551151)) @cascade { 
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
	index @filter(eq(ticker, "AAPL") or eq(ticker, "USD.FEDFUND")) {
		timestamp
		ticker
		value
	}
}`

	actual := GetImpliedVolRequestQuery("1513551151", "AAPL", "USD.FEDFUND")
	assert.EqualValues(t, expected, actual)
}

func BenchmarkGetImpliedVolRequestQuery(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetImpliedVolRequestQuery("1513551151", "AAPL", "USD.FEDFUND")
	}
}
