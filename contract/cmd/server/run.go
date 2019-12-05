package main

import (
	"github.com/gooption-io/gooption/v1/contract/app/server"

	"github.com/srvc/appctx"

	"github.com/izumin5210/grapi/pkg/grapiserver"
)

func run() error {
	// Application context
	ctx := appctx.Global()

	s := grapiserver.New(
		grapiserver.WithDefaultLogger(),
		grapiserver.WithServers(
			server.NewEuropeanServiceServer(),
			server.NewStockQuoteServiceServer(),
		),
	)
	return s.Serve(ctx)
}
