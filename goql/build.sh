# protoc \
#     --proto_path=. \
#     --proto_path=$GOPATH/src \
#     --proto_path=$GOPATH/src/github.com/gooption-io/gooption/proto \
#     --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
#     --cpp_out=. \
#     service.proto

# protoc \
#     --proto_path=. \
#     --proto_path=$GOPATH/src \
#     --proto_path=$GOPATH/src/github.com/gooption-io/gooption/proto \
#     --proto_path=$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
#     --grpc_out=. \
#     --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` \
#     service.proto

g++ \
    -std=c++14 -stdlib=libc++ \
    -O2 -Wall \
    -I /usr/local/include/boost \
    -I /usr/local/include/spdlog \
    -I . \
    -L/usr/local/lib \
    -lQuantLib \
    -lboost_program_options \
    -lgrpc++_reflection \
    -lprotobuf -lgrpc++ -lgrpc \
    *.cc \
    -o goql