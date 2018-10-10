protoc  \
	--proto_path=${GOPATH}/src \
	--proto_path=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
	--proto_path=. \
	--gogo_out=. \
	--govalidators_out=gogoimport=true:. \
	--weave_out=gogoimport=true:. \
	*.proto
