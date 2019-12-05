package main

import (
	"gobs/app/server"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"github.com/srvc/appctx"
)

func run() error {
	// Application context
	ctx := appctx.Global()

	s := grapiserver.New(
		grapiserver.WithDefaultLogger(),
		grapiserver.WithServers(
			server.NewEuropeanPricerServiceServer(),
		),
	)
	return s.Serve(ctx)
}
