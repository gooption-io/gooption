package server

import (
	"context"
	"encoding/json"
	"fmt"
	api_contract_pb "github.com/gooption-io/gooption/contract/api"
	api_pb "github.com/gooption-io/gooption/v1/gobs/api"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_europeanPricerServiceServerImpl_Compute(t *testing.T) {
	svr := NewEuropeanPricerServiceServer()
	ctx := context.Background()

	price_req := &api_pb.ComputationRequest{
		Pricingdate: 1564669224,

		// european option
		Contract: &api_contract_pb.European{
			Strike:    159.76,
			Putcall:   "put",
			Undticker: "AAPL",
			Expiry:    1582890485,
		},

		// market data
		Vol:  0.1,
		Rate: 0.01,
		Spot: &api_contract_pb.StockQuote{
			Close: 264.16,
		},
	}

	resp, err := svr.Compute(ctx, price_req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotEmpty(t, resp.Price)

	json_req, _ := json.MarshalIndent(price_req, "", "\t")
	fmt.Printf("%s \n", json_req)

	json_resp, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Printf("%s \n", json_resp)
}
