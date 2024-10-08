// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: contract.proto

package pb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import encoding_binary "encoding/binary"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type European struct {
	Timestamp            float64  `protobuf:"fixed64,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Ticker               string   `protobuf:"bytes,2,opt,name=ticker,proto3" json:"ticker,omitempty"`
	Undticker            string   `protobuf:"bytes,3,opt,name=undticker,proto3" json:"undticker,omitempty"`
	Strike               float64  `protobuf:"fixed64,4,opt,name=strike,proto3" json:"strike,omitempty"`
	Expiry               float64  `protobuf:"fixed64,5,opt,name=expiry,proto3" json:"expiry,omitempty"`
	Putcall              string   `protobuf:"bytes,6,opt,name=putcall,proto3" json:"putcall,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *European) Reset()         { *m = European{} }
func (m *European) String() string { return proto.CompactTextString(m) }
func (*European) ProtoMessage()    {}
func (*European) Descriptor() ([]byte, []int) {
	return fileDescriptor_contract_321c0fd306093d94, []int{0}
}
func (m *European) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *European) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_European.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (dst *European) XXX_Merge(src proto.Message) {
	xxx_messageInfo_European.Merge(dst, src)
}
func (m *European) XXX_Size() int {
	return m.Size()
}
func (m *European) XXX_DiscardUnknown() {
	xxx_messageInfo_European.DiscardUnknown(m)
}

var xxx_messageInfo_European proto.InternalMessageInfo

func (m *European) GetTimestamp() float64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *European) GetTicker() string {
	if m != nil {
		return m.Ticker
	}
	return ""
}

func (m *European) GetUndticker() string {
	if m != nil {
		return m.Undticker
	}
	return ""
}

func (m *European) GetStrike() float64 {
	if m != nil {
		return m.Strike
	}
	return 0
}

func (m *European) GetExpiry() float64 {
	if m != nil {
		return m.Expiry
	}
	return 0
}

func (m *European) GetPutcall() string {
	if m != nil {
		return m.Putcall
	}
	return ""
}

func init() {
	proto.RegisterType((*European)(nil), "pb.European")
}
func (m *European) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *European) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Timestamp != 0 {
		dAtA[i] = 0x9
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Timestamp))))
		i += 8
	}
	if len(m.Ticker) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Ticker)))
		i += copy(dAtA[i:], m.Ticker)
	}
	if len(m.Undticker) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Undticker)))
		i += copy(dAtA[i:], m.Undticker)
	}
	if m.Strike != 0 {
		dAtA[i] = 0x21
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Strike))))
		i += 8
	}
	if m.Expiry != 0 {
		dAtA[i] = 0x29
		i++
		encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.Expiry))))
		i += 8
	}
	if len(m.Putcall) > 0 {
		dAtA[i] = 0x32
		i++
		i = encodeVarintContract(dAtA, i, uint64(len(m.Putcall)))
		i += copy(dAtA[i:], m.Putcall)
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintContract(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *European) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Timestamp != 0 {
		n += 9
	}
	l = len(m.Ticker)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	l = len(m.Undticker)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	if m.Strike != 0 {
		n += 9
	}
	if m.Expiry != 0 {
		n += 9
	}
	l = len(m.Putcall)
	if l > 0 {
		n += 1 + l + sovContract(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovContract(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozContract(x uint64) (n int) {
	return sovContract(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *European) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowContract
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: European: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: European: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Timestamp = float64(math.Float64frombits(v))
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ticker", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ticker = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Undticker", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Undticker = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Strike", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Strike = float64(math.Float64frombits(v))
		case 5:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field Expiry", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.Expiry = float64(math.Float64frombits(v))
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Putcall", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowContract
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthContract
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Putcall = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipContract(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthContract
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipContract(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowContract
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowContract
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowContract
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthContract
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowContract
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipContract(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthContract = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowContract   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("contract.proto", fileDescriptor_contract_321c0fd306093d94) }

var fileDescriptor_contract_321c0fd306093d94 = []byte{
	// 174 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0xce, 0xcf, 0x2b,
	0x29, 0x4a, 0x4c, 0x2e, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x5a,
	0xc1, 0xc8, 0xc5, 0xe1, 0x5a, 0x5a, 0x94, 0x5f, 0x90, 0x9a, 0x98, 0x27, 0x24, 0xc3, 0xc5, 0x59,
	0x92, 0x99, 0x9b, 0x5a, 0x5c, 0x92, 0x98, 0x5b, 0x20, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x18, 0x84,
	0x10, 0x10, 0x12, 0xe3, 0x62, 0x2b, 0xc9, 0x4c, 0xce, 0x4e, 0x2d, 0x92, 0x60, 0x52, 0x60, 0xd4,
	0xe0, 0x0c, 0x82, 0xf2, 0x40, 0xba, 0x4a, 0xf3, 0x52, 0xa0, 0x52, 0xcc, 0x60, 0x29, 0x84, 0x00,
	0x48, 0x57, 0x71, 0x49, 0x51, 0x66, 0x76, 0xaa, 0x04, 0x0b, 0xd8, 0x40, 0x28, 0x0f, 0x24, 0x9e,
	0x5a, 0x51, 0x90, 0x59, 0x54, 0x29, 0xc1, 0x0a, 0x11, 0x87, 0xf0, 0x84, 0x24, 0xb8, 0xd8, 0x0b,
	0x4a, 0x4b, 0x92, 0x13, 0x73, 0x72, 0x24, 0xd8, 0xc0, 0x66, 0xc1, 0xb8, 0x4e, 0x3c, 0x27, 0x1e,
	0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0x63, 0x12, 0x1b, 0xd8, 0x0f, 0xc6,
	0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfc, 0x5d, 0x99, 0x1e, 0xd5, 0x00, 0x00, 0x00,
}
