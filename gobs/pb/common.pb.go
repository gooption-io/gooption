// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: common.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type OptionType int32

const (
	OptionType_PUT  OptionType = 0
	OptionType_CALL OptionType = 1
)

var OptionType_name = map[int32]string{
	0: "PUT",
	1: "CALL",
}
var OptionType_value = map[string]int32{
	"PUT":  0,
	"CALL": 1,
}

func (x OptionType) String() string {
	return proto.EnumName(OptionType_name, int32(x))
}
func (OptionType) EnumDescriptor() ([]byte, []int) { return fileDescriptorCommon, []int{0} }

func init() {
	proto.RegisterEnum("pb.OptionType", OptionType_name, OptionType_value)
}

func init() { proto.RegisterFile("common.proto", fileDescriptorCommon) }

var fileDescriptorCommon = []byte{
	// 99 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0xcd,
	0xcd, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0xd2, 0x92, 0xe7, 0xe2,
	0xf2, 0x2f, 0x28, 0xc9, 0xcc, 0xcf, 0x0b, 0xa9, 0x2c, 0x48, 0x15, 0x62, 0xe7, 0x62, 0x0e, 0x08,
	0x0d, 0x11, 0x60, 0x10, 0xe2, 0xe0, 0x62, 0x71, 0x76, 0xf4, 0xf1, 0x11, 0x60, 0x74, 0x12, 0x38,
	0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23, 0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x67, 0x3c, 0x96, 0x63,
	0x48, 0x62, 0x03, 0xeb, 0x36, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xd8, 0x39, 0x22, 0x86, 0x4d,
	0x00, 0x00, 0x00,
}
