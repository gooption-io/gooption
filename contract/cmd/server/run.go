package main

import (
	"contract/app/server"

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
		),
	)
	return s.Serve(ctx)
}
