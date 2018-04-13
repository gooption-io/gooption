package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPriceRequestQuery(t *testing.T) {
	expected := `
{
	contract(func: eq(ticker, "AAPL DEC2017 PUT")){
		ticker
		strike
		und as undticker
		expiry
		putcall
	}

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
}

fragment indexInfo {
	index @filter(eq(ticker, val(und)) or eq(ticker, "USD.FEDFUND")) {
		timestamp
		ticker
		value
	}
}`

	actual := GetPriceRequestQuery("1513551151", "AAPL DEC2017 PUT", "USD.FEDFUND")
	assert.EqualValues(t, expected, actual)
}

func BenchmarkGetPriceRequestQuery(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetPriceRequestQuery("1513551151", "AAPL DEC2017 PUT", "USD.FEDFUND")
	}
}
