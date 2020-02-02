module github.com/gooption-io/gooption/v1/gobs

go 1.13

replace github.com/gooption-io/gooption/contract => /Users/lehajam/go/src/github.com/gooption-io/gooption/contract

require (
	github.com/golang/protobuf v1.3.2
	github.com/gooption-io/gooption/contract v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/izumin5210/gex v0.5.1
	github.com/izumin5210/grapi v0.5.0
	github.com/srvc/appctx v0.1.0
	github.com/stretchr/testify v1.4.0
	gonum.org/v1/gonum v0.0.0-20191004082826-a11ea52b6f3c
	google.golang.org/genproto v0.0.0-20200127141224-2548664c049f
	google.golang.org/grpc v1.26.0
)
