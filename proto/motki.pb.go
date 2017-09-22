// Code generated by protoc-gen-go. DO NOT EDIT.
// source: motki.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	motki.proto
	model.proto
	evedb.proto

It has these top-level messages:
	Result
	Token
	AuthenticateRequest
	AuthenticateResponse
	Character
	GetCharacterRequest
	CharacterResponse
	Corporation
	GetCorporationRequest
	CorporationResponse
	Alliance
	GetAllianceRequest
	AllianceResponse
	Product
	ProductResponse
	GetProductRequest
	NewProductRequest
	SaveProductRequest
	GetProductsRequest
	ProductsResponse
	MarketPrice
	GetMarketPriceRequest
	GetMarketPriceResponse
	Blueprint
	GetCorpBlueprintsRequest
	GetCorpBlueprintsResponse
	Icon
	Race
	Ancestry
	Bloodline
	System
	Constellation
	Region
	ItemType
	ItemTypeDetail
	MaterialSheet
	Material
	GetRegionRequest
	GetRegionResponse
	GetRegionsRequest
	GetRegionsResponse
	GetConstellationRequest
	GetConstellationResponse
	GetSystemRequest
	GetSystemResponse
	GetRaceRequest
	GetRaceResponse
	GetRacesRequest
	GetRacesResponse
	GetBloodlineRequest
	GetBloodlineResponse
	GetAncestryRequest
	GetAncestryResponse
	GetItemTypeRequest
	GetItemTypeResponse
	GetItemTypeDetailRequest
	GetItemTypeDetailResponse
	QueryItemTypesRequest
	QueryItemTypesResponse
	QueryItemTypeDetailsRequest
	QueryItemTypeDetailsResponse
	GetMaterialSheetRequest
	GetMaterialSheetResponse
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// A Status indicates success or failure.
type Status int32

const (
	Status_FAILURE Status = 0
	Status_SUCCESS Status = 1
)

var Status_name = map[int32]string{
	0: "FAILURE",
	1: "SUCCESS",
}
var Status_value = map[string]int32{
	"FAILURE": 0,
	"SUCCESS": 1,
}

func (x Status) String() string {
	return proto1.EnumName(Status_name, int32(x))
}
func (Status) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// A Role describes a user's authorization for certain resources within the system.
type Role int32

const (
	Role_ANON      Role = 0
	Role_USER      Role = 1
	Role_LOGISTICS Role = 2
)

var Role_name = map[int32]string{
	0: "ANON",
	1: "USER",
	2: "LOGISTICS",
}
var Role_value = map[string]int32{
	"ANON":      0,
	"USER":      1,
	"LOGISTICS": 2,
}

func (x Role) String() string {
	return proto1.EnumName(Role_name, int32(x))
}
func (Role) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// A Result contains a status and an optional description.
type Result struct {
	Status Status `protobuf:"varint,1,opt,name=status,enum=motki.Status" json:"status,omitempty"`
	// Description contains some text about a failure in most cases.
	Description string `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto1.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Result) GetStatus() Status {
	if m != nil {
		return m.Status
	}
	return Status_FAILURE
}

func (m *Result) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

// A Token contains a session identifier representing a valid user session.
type Token struct {
	Identifier string `protobuf:"bytes,1,opt,name=identifier" json:"identifier,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto1.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Token) GetIdentifier() string {
	if m != nil {
		return m.Identifier
	}
	return ""
}

type AuthenticateRequest struct {
	Username string `protobuf:"bytes,1,opt,name=username" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *AuthenticateRequest) Reset()                    { *m = AuthenticateRequest{} }
func (m *AuthenticateRequest) String() string            { return proto1.CompactTextString(m) }
func (*AuthenticateRequest) ProtoMessage()               {}
func (*AuthenticateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AuthenticateRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *AuthenticateRequest) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type AuthenticateResponse struct {
	Result *Result `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	Token  *Token  `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
}

func (m *AuthenticateResponse) Reset()                    { *m = AuthenticateResponse{} }
func (m *AuthenticateResponse) String() string            { return proto1.CompactTextString(m) }
func (*AuthenticateResponse) ProtoMessage()               {}
func (*AuthenticateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AuthenticateResponse) GetResult() *Result {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *AuthenticateResponse) GetToken() *Token {
	if m != nil {
		return m.Token
	}
	return nil
}

func init() {
	proto1.RegisterType((*Result)(nil), "motki.Result")
	proto1.RegisterType((*Token)(nil), "motki.Token")
	proto1.RegisterType((*AuthenticateRequest)(nil), "motki.AuthenticateRequest")
	proto1.RegisterType((*AuthenticateResponse)(nil), "motki.AuthenticateResponse")
	proto1.RegisterEnum("motki.Status", Status_name, Status_value)
	proto1.RegisterEnum("motki.Role", Role_name, Role_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for AuthenticationService service

type AuthenticationServiceClient interface {
	Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error)
}

type authenticationServiceClient struct {
	cc *grpc.ClientConn
}

func NewAuthenticationServiceClient(cc *grpc.ClientConn) AuthenticationServiceClient {
	return &authenticationServiceClient{cc}
}

func (c *authenticationServiceClient) Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	out := new(AuthenticateResponse)
	err := grpc.Invoke(ctx, "/motki.AuthenticationService/Authenticate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AuthenticationService service

type AuthenticationServiceServer interface {
	Authenticate(context.Context, *AuthenticateRequest) (*AuthenticateResponse, error)
}

func RegisterAuthenticationServiceServer(s *grpc.Server, srv AuthenticationServiceServer) {
	s.RegisterService(&_AuthenticationService_serviceDesc, srv)
}

func _AuthenticationService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthenticationServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/motki.AuthenticationService/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthenticationServiceServer).Authenticate(ctx, req.(*AuthenticateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AuthenticationService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "motki.AuthenticationService",
	HandlerType: (*AuthenticationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authenticate",
			Handler:    _AuthenticationService_Authenticate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "motki.proto",
}

func init() { proto1.RegisterFile("motki.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 325 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xd1, 0x6b, 0xc2, 0x30,
	0x10, 0xc6, 0xad, 0xd8, 0xaa, 0x57, 0x1d, 0x25, 0xdb, 0x40, 0x1c, 0x0c, 0x29, 0x8c, 0x39, 0x1f,
	0x7c, 0x70, 0x7f, 0x81, 0x13, 0x27, 0x82, 0x53, 0x96, 0xe8, 0xcb, 0x9e, 0xd6, 0xe9, 0x8d, 0x05,
	0xb5, 0xe9, 0x92, 0x74, 0xfb, 0xf7, 0x47, 0xd2, 0x20, 0x1d, 0xf8, 0xd4, 0x7e, 0xf7, 0x5d, 0xbe,
	0xdc, 0xef, 0x02, 0xe1, 0x51, 0xe8, 0x3d, 0x1f, 0x66, 0x52, 0x68, 0x41, 0x7c, 0x2b, 0xe2, 0x57,
	0x08, 0x28, 0xaa, 0xfc, 0xa0, 0xc9, 0x1d, 0x04, 0x4a, 0x27, 0x3a, 0x57, 0x1d, 0xaf, 0xe7, 0xf5,
	0x2f, 0x46, 0xed, 0x61, 0xd1, 0xce, 0x6c, 0x91, 0x3a, 0x93, 0xf4, 0x20, 0xdc, 0xa1, 0xda, 0x4a,
	0x9e, 0x69, 0x2e, 0xd2, 0x4e, 0xb5, 0xe7, 0xf5, 0x9b, 0xb4, 0x5c, 0x8a, 0xef, 0xc1, 0x5f, 0x8b,
	0x3d, 0xa6, 0xe4, 0x16, 0x80, 0xef, 0x30, 0xd5, 0xfc, 0x93, 0xa3, 0xb4, 0xa9, 0x4d, 0x5a, 0xaa,
	0xc4, 0x2f, 0x70, 0x39, 0xce, 0xf5, 0x97, 0xd1, 0xdb, 0x44, 0x23, 0xc5, 0xef, 0x1c, 0x95, 0x26,
	0x5d, 0x68, 0xe4, 0x0a, 0x65, 0x9a, 0x1c, 0xd1, 0x1d, 0x3a, 0x69, 0xe3, 0x65, 0x89, 0x52, 0xbf,
	0x42, 0xee, 0xdc, 0xd5, 0x27, 0x1d, 0x27, 0x70, 0xf5, 0x3f, 0x4e, 0x65, 0x22, 0x55, 0x68, 0xc0,
	0xa4, 0x45, 0xb4, 0x69, 0xe1, 0x09, 0xac, 0xe0, 0xa6, 0xce, 0x24, 0x31, 0xf8, 0xda, 0x8c, 0x6d,
	0x73, 0xc3, 0x51, 0xcb, 0x75, 0x59, 0x14, 0x5a, 0x58, 0x83, 0x18, 0x82, 0x62, 0x1d, 0x24, 0x84,
	0xfa, 0xf3, 0x78, 0xbe, 0xd8, 0xd0, 0x69, 0x54, 0x31, 0x82, 0x6d, 0x26, 0x93, 0x29, 0x63, 0x91,
	0x37, 0x78, 0x80, 0x1a, 0x15, 0x07, 0x24, 0x0d, 0xa8, 0x8d, 0x97, 0xab, 0x65, 0x54, 0x31, 0x7f,
	0x1b, 0x36, 0xa5, 0x91, 0x47, 0xda, 0xd0, 0x5c, 0xac, 0x66, 0x73, 0xb6, 0x9e, 0x4f, 0x58, 0x54,
	0x1d, 0xbd, 0xc3, 0x75, 0x69, 0x62, 0x2e, 0x52, 0x86, 0xf2, 0x87, 0x6f, 0x91, 0xcc, 0xa0, 0x55,
	0x46, 0x21, 0x5d, 0x37, 0xcc, 0x99, 0x75, 0x75, 0x6f, 0xce, 0x7a, 0x05, 0xfb, 0x53, 0xfd, 0xcd,
	0xb7, 0xcf, 0xfd, 0x11, 0xd8, 0xcf, 0xe3, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x60, 0xfa, 0x4f,
	0xea, 0x04, 0x02, 0x00, 0x00,
}
