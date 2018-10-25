protoc \
	--proto_path=$GOPATH/src \
	--proto_path=$GOPATH/src/github.com/gogo/protobuf/protobuf \
	--proto_path=. \
	--gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. \
	*.proto
