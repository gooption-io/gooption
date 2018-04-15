package query

import (
	"github.com/valyala/fasttemplate"
)

var (

	// Query with optional tags goes here
	impliedvolrequest = `
{
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

	quotes(func: eq(timestamp, {{timestamp}})) @cascade { 
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
	index @filter(eq(ticker, "{{undTicker}}") or eq(ticker, "{{rateTicker}}")) {
		timestamp
		ticker
		value
	}
}`

	impliedvolrequestTemplate = fasttemplate.New(impliedvolrequest, startTag, endTag)
)

func init() {
	// For template reflection
	AllTemplates["impliedvolrequest"] = impliedvolrequestTemplate
}

func GetImpliedVolRequestQuery(timestamp, undTicker, rateTicker string) string {

	return impliedvolrequestTemplate.ExecuteString(map[string]interface{}{
		"timestamp": timestamp, "undTicker": undTicker, "rateTicker": rateTicker,
	})
}
