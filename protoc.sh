protoc  \
    -I contract/api/protos -I gobs/api/protos  \
    -I /Users/lehajam/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.8.5  \
    -I /Users/lehajam/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.8.5/third_party/googleapis  \
    --go_out=plugins=grpc,paths=source_relative:gobs/api gobs/api/protos/european_pricer.proto