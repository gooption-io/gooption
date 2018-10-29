protoc \
    --proto_path=. \
    --proto_path=$GOPATH/src \
    --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --gogofast_out=plugins=grpc:go/pb \
    --cpp_out=cpp \
    *.proto

protoc \
    --proto_path=. \
    --proto_path=$GOPATH/src \
    --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --grpc_out=cpp \
    --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` \
    service.proto

protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true,grpc_api_configuration=service.yaml:go/pb \
  service.proto
