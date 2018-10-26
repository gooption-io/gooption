protoc \
    --proto_path=. \
    --proto_path=$GOPATH/src \
    --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --gogofast_out=plugins=grpc:. \
    --grpc-gateway_out=logtostderr=true:. \
    *.proto
