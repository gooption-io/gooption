# protoc \
#     --proto_path=. \
#     --proto_path=$GOPATH/src \
#     --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
#     --gogofast_out=plugins=grpc:go/pb \
#     --grpc-gateway_out=logtostderr=true:go/pb \
#     --cpp_out=cpp \
#     --grpc_out=cpp \
#     --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` \
#     *.proto

protoc \
    --proto_path=. \
    --proto_path=$GOPATH/src \
    --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    --cpp_out=cpp \
    --grpc_out=cpp \
    --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` \
    *.proto
