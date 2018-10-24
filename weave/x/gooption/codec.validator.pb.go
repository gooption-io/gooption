package gooption

import (
	fmt "fmt"

	proto "github.com/gogo/protobuf/proto"
	"github.com/iov-one/weave"

	math "math"

	_ "github.com/gogo/protobuf/gogoproto"

	_ "github.com/lehajam/protoc-gen-weave/x/bucket"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

func (this *Valuation) Validate() error {
	return weave.Address(this.Sender).Validate()
}
func (this *CreateValuationMsg) Validate() error {
	return weave.Address(this.Sender).Validate()
}
