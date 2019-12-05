// Code generated by protoc-gen-go. DO NOT EDIT.
// source: european_pricer.proto

package api_pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	api "github.com/gooption-io/gooption/contract/api"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ComputationRequest struct {
	Pricingdate          float64         `protobuf:"fixed64,2,opt,name=pricingdate,proto3" json:"pricingdate,omitempty"`
	Contract             *api.European   `protobuf:"bytes,3,opt,name=contract,proto3" json:"contract,omitempty"`
	Spot                 *api.StockQuote `protobuf:"bytes,4,opt,name=spot,proto3" json:"spot,omitempty"`
	Vol                  float64         `protobuf:"fixed64,5,opt,name=vol,proto3" json:"vol,omitempty"`
	Rate                 float64         `protobuf:"fixed64,6,opt,name=rate,proto3" json:"rate,omitempty"`
	Greek                []string        `protobuf:"bytes,7,rep,name=greek,proto3" json:"greek,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *ComputationRequest) Reset()         { *m = ComputationRequest{} }
func (m *ComputationRequest) String() string { return proto.CompactTextString(m) }
func (*ComputationRequest) ProtoMessage()    {}
func (*ComputationRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd5af2df7a4d0db1, []int{0}
}

func (m *ComputationRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComputationRequest.Unmarshal(m, b)
}
func (m *ComputationRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComputationRequest.Marshal(b, m, deterministic)
}
func (m *ComputationRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComputationRequest.Merge(m, src)
}
func (m *ComputationRequest) XXX_Size() int {
	return xxx_messageInfo_ComputationRequest.Size(m)
}
func (m *ComputationRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ComputationRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ComputationRequest proto.InternalMessageInfo

func (m *ComputationRequest) GetPricingdate() float64 {
	if m != nil {
		return m.Pricingdate
	}
	return 0
}

func (m *ComputationRequest) GetContract() *api.European {
	if m != nil {
		return m.Contract
	}
	return nil
}

func (m *ComputationRequest) GetSpot() *api.StockQuote {
	if m != nil {
		return m.Spot
	}
	return nil
}

func (m *ComputationRequest) GetVol() float64 {
	if m != nil {
		return m.Vol
	}
	return 0
}

func (m *ComputationRequest) GetRate() float64 {
	if m != nil {
		return m.Rate
	}
	return 0
}

func (m *ComputationRequest) GetGreek() []string {
	if m != nil {
		return m.Greek
	}
	return nil
}

type ComputationResponse struct {
	Price                float64  `protobuf:"fixed64,2,opt,name=price,proto3" json:"price,omitempty"`
	Greeks               []*Greek `protobuf:"bytes,3,rep,name=greeks,proto3" json:"greeks,omitempty"`
	PriceError           string   `protobuf:"bytes,4,opt,name=price_error,json=priceError,proto3" json:"price_error,omitempty"`
	GreeksError          string   `protobuf:"bytes,5,opt,name=greeks_error,json=greeksError,proto3" json:"greeks_error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ComputationResponse) Reset()         { *m = ComputationResponse{} }
func (m *ComputationResponse) String() string { return proto.CompactTextString(m) }
func (*ComputationResponse) ProtoMessage()    {}
func (*ComputationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd5af2df7a4d0db1, []int{1}
}

func (m *ComputationResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ComputationResponse.Unmarshal(m, b)
}
func (m *ComputationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ComputationResponse.Marshal(b, m, deterministic)
}
func (m *ComputationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComputationResponse.Merge(m, src)
}
func (m *ComputationResponse) XXX_Size() int {
	return xxx_messageInfo_ComputationResponse.Size(m)
}
func (m *ComputationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ComputationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ComputationResponse proto.InternalMessageInfo

func (m *ComputationResponse) GetPrice() float64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *ComputationResponse) GetGreeks() []*Greek {
	if m != nil {
		return m.Greeks
	}
	return nil
}

func (m *ComputationResponse) GetPriceError() string {
	if m != nil {
		return m.PriceError
	}
	return ""
}

func (m *ComputationResponse) GetGreeksError() string {
	if m != nil {
		return m.GreeksError
	}
	return ""
}

type Greek struct {
	Label                string   `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	Value                float64  `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
	Error                string   `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Greek) Reset()         { *m = Greek{} }
func (m *Greek) String() string { return proto.CompactTextString(m) }
func (*Greek) ProtoMessage()    {}
func (*Greek) Descriptor() ([]byte, []int) {
	return fileDescriptor_dd5af2df7a4d0db1, []int{2}
}

func (m *Greek) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Greek.Unmarshal(m, b)
}
func (m *Greek) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Greek.Marshal(b, m, deterministic)
}
func (m *Greek) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Greek.Merge(m, src)
}
func (m *Greek) XXX_Size() int {
	return xxx_messageInfo_Greek.Size(m)
}
func (m *Greek) XXX_DiscardUnknown() {
	xxx_messageInfo_Greek.DiscardUnknown(m)
}

var xxx_messageInfo_Greek proto.InternalMessageInfo

func (m *Greek) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *Greek) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *Greek) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func init() {
	proto.RegisterType((*ComputationRequest)(nil), "gooption.gobs.ComputationRequest")
	proto.RegisterType((*ComputationResponse)(nil), "gooption.gobs.ComputationResponse")
	proto.RegisterType((*Greek)(nil), "gooption.gobs.Greek")
}

func init() { proto.RegisterFile("european_pricer.proto", fileDescriptor_dd5af2df7a4d0db1) }

var fileDescriptor_dd5af2df7a4d0db1 = []byte{
	// 426 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xb1, 0x8e, 0xd4, 0x30,
	0x10, 0x86, 0x95, 0xcb, 0x66, 0xef, 0x6e, 0x02, 0xe8, 0x30, 0x77, 0x52, 0xb4, 0x80, 0xc8, 0xa5,
	0x4a, 0x01, 0x89, 0x6e, 0x29, 0x28, 0xe8, 0x40, 0x27, 0x44, 0x07, 0xd9, 0x8e, 0x26, 0x72, 0x82,
	0x15, 0xa2, 0xcd, 0x66, 0xbc, 0xb6, 0xb3, 0x0d, 0x1d, 0xaf, 0x40, 0xcd, 0x53, 0xd1, 0x53, 0xf1,
	0x20, 0xc8, 0xe3, 0x24, 0x62, 0x01, 0xd1, 0x79, 0x7e, 0x7f, 0x33, 0xf9, 0xe7, 0x8f, 0xe1, 0x4a,
	0x0c, 0x0a, 0xa5, 0xe0, 0x7d, 0x29, 0x55, 0x5b, 0x0b, 0x95, 0x49, 0x85, 0x06, 0xd9, 0xdd, 0x06,
	0x51, 0x9a, 0x16, 0xfb, 0xac, 0xc1, 0x4a, 0xaf, 0x1e, 0x35, 0x88, 0x4d, 0x27, 0x72, 0x2e, 0xdb,
	0x9c, 0xf7, 0x3d, 0x1a, 0x6e, 0xef, 0xb4, 0x83, 0x57, 0xf7, 0xa6, 0x19, 0x63, 0x7d, 0x5f, 0x1b,
	0xac, 0xb7, 0xe5, 0x7e, 0x40, 0x23, 0x9c, 0x94, 0xfc, 0xf0, 0x80, 0xbd, 0xc6, 0x9d, 0x1c, 0x5c,
	0x67, 0x21, 0xf6, 0x83, 0xd0, 0x86, 0xc5, 0x10, 0xda, 0xcf, 0xb6, 0x7d, 0xf3, 0x91, 0x1b, 0x11,
	0x9d, 0xc4, 0x5e, 0xea, 0x15, 0xbf, 0x4b, 0xec, 0x05, 0x9c, 0xd5, 0xd8, 0x1b, 0xc5, 0x6b, 0x13,
	0xf9, 0xb1, 0x97, 0x86, 0xeb, 0x87, 0xd9, 0xec, 0x6d, 0xba, 0xc9, 0x6e, 0x47, 0x03, 0xc5, 0x0c,
	0xb3, 0x1b, 0x58, 0x68, 0x89, 0x26, 0x5a, 0x50, 0xd3, 0xe3, 0x7f, 0x34, 0x6d, 0xac, 0xcb, 0xf7,
	0xd6, 0x64, 0x41, 0x28, 0xbb, 0x00, 0xff, 0x80, 0x5d, 0x14, 0x90, 0x0b, 0x7b, 0x64, 0x0c, 0x16,
	0xca, 0x1a, 0x5b, 0x92, 0x44, 0x67, 0x76, 0x09, 0x41, 0xa3, 0x84, 0xd8, 0x46, 0xa7, 0xb1, 0x9f,
	0x9e, 0x17, 0xae, 0x48, 0xbe, 0x79, 0xf0, 0xe0, 0x68, 0x41, 0x2d, 0xb1, 0xd7, 0x44, 0x53, 0xb0,
	0xe3, 0x6e, 0xae, 0x60, 0x4f, 0x61, 0x49, 0x6d, 0x3a, 0xf2, 0x63, 0x3f, 0x0d, 0xd7, 0x97, 0xd9,
	0x51, 0xde, 0xd9, 0x1b, 0x7b, 0x59, 0x8c, 0x0c, 0x7b, 0xe2, 0x52, 0x12, 0xa5, 0x50, 0x0a, 0x15,
	0x6d, 0x74, 0x5e, 0x00, 0x49, 0xb7, 0x56, 0x61, 0xd7, 0x70, 0xc7, 0xa1, 0x23, 0x11, 0x10, 0x11,
	0x3a, 0x8d, 0x90, 0xe4, 0x2d, 0x04, 0x34, 0xd4, 0x1a, 0xea, 0x78, 0x25, 0xba, 0xc8, 0x23, 0xc8,
	0x15, 0x56, 0x3d, 0xf0, 0x6e, 0x98, 0x6d, 0x52, 0x61, 0x55, 0x37, 0xd0, 0x77, 0x2c, 0x15, 0xeb,
	0xcf, 0x70, 0x35, 0xe5, 0xfd, 0x8e, 0xde, 0xcc, 0x46, 0xa8, 0x83, 0xdd, 0xaa, 0x82, 0x53, 0x17,
	0x81, 0x60, 0xd7, 0x7f, 0x2c, 0xf4, 0xf7, 0xbf, 0x5f, 0x25, 0xff, 0x43, 0x5c, 0x7a, 0xc9, 0xc5,
	0x97, 0xef, 0x3f, 0xbf, 0x9e, 0x00, 0x3b, 0xcb, 0x6b, 0x37, 0xf8, 0xd5, 0xcd, 0x87, 0xbc, 0x69,
	0xcd, 0xa7, 0xa1, 0xca, 0x6a, 0xdc, 0xe5, 0xd3, 0x84, 0x67, 0x2d, 0xce, 0xe7, 0xdc, 0x4e, 0xb3,
	0x2f, 0xf5, 0x25, 0x97, 0x6d, 0x29, 0xab, 0x6a, 0x49, 0x4f, 0xf0, 0xf9, 0xaf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x2b, 0xbc, 0x42, 0xab, 0xeb, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EuropeanPricerServiceClient is the client API for EuropeanPricerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EuropeanPricerServiceClient interface {
	Compute(ctx context.Context, in *ComputationRequest, opts ...grpc.CallOption) (*ComputationResponse, error)
}

type europeanPricerServiceClient struct {
	cc *grpc.ClientConn
}

func NewEuropeanPricerServiceClient(cc *grpc.ClientConn) EuropeanPricerServiceClient {
	return &europeanPricerServiceClient{cc}
}

func (c *europeanPricerServiceClient) Compute(ctx context.Context, in *ComputationRequest, opts ...grpc.CallOption) (*ComputationResponse, error) {
	out := new(ComputationResponse)
	err := c.cc.Invoke(ctx, "/gooption.gobs.EuropeanPricerService/Compute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EuropeanPricerServiceServer is the server API for EuropeanPricerService service.
type EuropeanPricerServiceServer interface {
	Compute(context.Context, *ComputationRequest) (*ComputationResponse, error)
}

// UnimplementedEuropeanPricerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedEuropeanPricerServiceServer struct {
}

func (*UnimplementedEuropeanPricerServiceServer) Compute(ctx context.Context, req *ComputationRequest) (*ComputationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Compute not implemented")
}

func RegisterEuropeanPricerServiceServer(s *grpc.Server, srv EuropeanPricerServiceServer) {
	s.RegisterService(&_EuropeanPricerService_serviceDesc, srv)
}

func _EuropeanPricerService_Compute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ComputationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EuropeanPricerServiceServer).Compute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gooption.gobs.EuropeanPricerService/Compute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EuropeanPricerServiceServer).Compute(ctx, req.(*ComputationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EuropeanPricerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gooption.gobs.EuropeanPricerService",
	HandlerType: (*EuropeanPricerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Compute",
			Handler:    _EuropeanPricerService_Compute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "european_pricer.proto",
}
