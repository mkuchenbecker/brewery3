// Code generated by protoc-gen-go. DO NOT EDIT.
// source: worldmodel.proto

package deltav_model

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type IdentificationProperty struct {
	Transponder          string   `protobuf:"bytes,1,opt,name=transponder,proto3" json:"transponder,omitempty"`
	InternalId           string   `protobuf:"bytes,2,opt,name=internal_id,json=internalId,proto3" json:"internal_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IdentificationProperty) Reset()         { *m = IdentificationProperty{} }
func (m *IdentificationProperty) String() string { return proto.CompactTextString(m) }
func (*IdentificationProperty) ProtoMessage()    {}
func (*IdentificationProperty) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{0}
}

func (m *IdentificationProperty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IdentificationProperty.Unmarshal(m, b)
}
func (m *IdentificationProperty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IdentificationProperty.Marshal(b, m, deterministic)
}
func (m *IdentificationProperty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IdentificationProperty.Merge(m, src)
}
func (m *IdentificationProperty) XXX_Size() int {
	return xxx_messageInfo_IdentificationProperty.Size(m)
}
func (m *IdentificationProperty) XXX_DiscardUnknown() {
	xxx_messageInfo_IdentificationProperty.DiscardUnknown(m)
}

var xxx_messageInfo_IdentificationProperty proto.InternalMessageInfo

func (m *IdentificationProperty) GetTransponder() string {
	if m != nil {
		return m.Transponder
	}
	return ""
}

func (m *IdentificationProperty) GetInternalId() string {
	if m != nil {
		return m.InternalId
	}
	return ""
}

type RadiationProperty struct {
	Type                 RadiationType `protobuf:"varint,1,opt,name=type,proto3,enum=deltav.model.RadiationType" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *RadiationProperty) Reset()         { *m = RadiationProperty{} }
func (m *RadiationProperty) String() string { return proto.CompactTextString(m) }
func (*RadiationProperty) ProtoMessage()    {}
func (*RadiationProperty) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{1}
}

func (m *RadiationProperty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RadiationProperty.Unmarshal(m, b)
}
func (m *RadiationProperty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RadiationProperty.Marshal(b, m, deterministic)
}
func (m *RadiationProperty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RadiationProperty.Merge(m, src)
}
func (m *RadiationProperty) XXX_Size() int {
	return xxx_messageInfo_RadiationProperty.Size(m)
}
func (m *RadiationProperty) XXX_DiscardUnknown() {
	xxx_messageInfo_RadiationProperty.DiscardUnknown(m)
}

var xxx_messageInfo_RadiationProperty proto.InternalMessageInfo

func (m *RadiationProperty) GetType() RadiationType {
	if m != nil {
		return m.Type
	}
	return RadiationType_UNKNOWN_RADIATION_TYPE
}

type SensorProperty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SensorProperty) Reset()         { *m = SensorProperty{} }
func (m *SensorProperty) String() string { return proto.CompactTextString(m) }
func (*SensorProperty) ProtoMessage()    {}
func (*SensorProperty) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{2}
}

func (m *SensorProperty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SensorProperty.Unmarshal(m, b)
}
func (m *SensorProperty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SensorProperty.Marshal(b, m, deterministic)
}
func (m *SensorProperty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SensorProperty.Merge(m, src)
}
func (m *SensorProperty) XXX_Size() int {
	return xxx_messageInfo_SensorProperty.Size(m)
}
func (m *SensorProperty) XXX_DiscardUnknown() {
	xxx_messageInfo_SensorProperty.DiscardUnknown(m)
}

var xxx_messageInfo_SensorProperty proto.InternalMessageInfo

type DetectableProperty struct {
	// Types that are valid to be assigned to Property:
	//	*DetectableProperty_Id
	//	*DetectableProperty_Radiation
	Property             isDetectableProperty_Property `protobuf_oneof:"property"`
	Intensity            float32                       `protobuf:"fixed32,100,opt,name=intensity,proto3" json:"intensity,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *DetectableProperty) Reset()         { *m = DetectableProperty{} }
func (m *DetectableProperty) String() string { return proto.CompactTextString(m) }
func (*DetectableProperty) ProtoMessage()    {}
func (*DetectableProperty) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{3}
}

func (m *DetectableProperty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DetectableProperty.Unmarshal(m, b)
}
func (m *DetectableProperty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DetectableProperty.Marshal(b, m, deterministic)
}
func (m *DetectableProperty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DetectableProperty.Merge(m, src)
}
func (m *DetectableProperty) XXX_Size() int {
	return xxx_messageInfo_DetectableProperty.Size(m)
}
func (m *DetectableProperty) XXX_DiscardUnknown() {
	xxx_messageInfo_DetectableProperty.DiscardUnknown(m)
}

var xxx_messageInfo_DetectableProperty proto.InternalMessageInfo

type isDetectableProperty_Property interface {
	isDetectableProperty_Property()
}

type DetectableProperty_Id struct {
	Id *IdentificationProperty `protobuf:"bytes,1,opt,name=id,proto3,oneof"`
}

type DetectableProperty_Radiation struct {
	Radiation *RadiationProperty `protobuf:"bytes,2,opt,name=radiation,proto3,oneof"`
}

func (*DetectableProperty_Id) isDetectableProperty_Property() {}

func (*DetectableProperty_Radiation) isDetectableProperty_Property() {}

func (m *DetectableProperty) GetProperty() isDetectableProperty_Property {
	if m != nil {
		return m.Property
	}
	return nil
}

func (m *DetectableProperty) GetId() *IdentificationProperty {
	if x, ok := m.GetProperty().(*DetectableProperty_Id); ok {
		return x.Id
	}
	return nil
}

func (m *DetectableProperty) GetRadiation() *RadiationProperty {
	if x, ok := m.GetProperty().(*DetectableProperty_Radiation); ok {
		return x.Radiation
	}
	return nil
}

func (m *DetectableProperty) GetIntensity() float32 {
	if m != nil {
		return m.Intensity
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*DetectableProperty) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _DetectableProperty_OneofMarshaler, _DetectableProperty_OneofUnmarshaler, _DetectableProperty_OneofSizer, []interface{}{
		(*DetectableProperty_Id)(nil),
		(*DetectableProperty_Radiation)(nil),
	}
}

func _DetectableProperty_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*DetectableProperty)
	// property
	switch x := m.Property.(type) {
	case *DetectableProperty_Id:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Id); err != nil {
			return err
		}
	case *DetectableProperty_Radiation:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Radiation); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("DetectableProperty.Property has unexpected type %T", x)
	}
	return nil
}

func _DetectableProperty_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*DetectableProperty)
	switch tag {
	case 1: // property.id
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(IdentificationProperty)
		err := b.DecodeMessage(msg)
		m.Property = &DetectableProperty_Id{msg}
		return true, err
	case 2: // property.radiation
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RadiationProperty)
		err := b.DecodeMessage(msg)
		m.Property = &DetectableProperty_Radiation{msg}
		return true, err
	default:
		return false, nil
	}
}

func _DetectableProperty_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*DetectableProperty)
	// property
	switch x := m.Property.(type) {
	case *DetectableProperty_Id:
		s := proto.Size(x.Id)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *DetectableProperty_Radiation:
		s := proto.Size(x.Radiation)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type RegisterRequest struct {
	Id                   string                `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Position             *Position             `protobuf:"bytes,2,opt,name=position,proto3" json:"position,omitempty"`
	Properties           []*DetectableProperty `protobuf:"bytes,3,rep,name=properties,proto3" json:"properties,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *RegisterRequest) Reset()         { *m = RegisterRequest{} }
func (m *RegisterRequest) String() string { return proto.CompactTextString(m) }
func (*RegisterRequest) ProtoMessage()    {}
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{4}
}

func (m *RegisterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterRequest.Unmarshal(m, b)
}
func (m *RegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterRequest.Marshal(b, m, deterministic)
}
func (m *RegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterRequest.Merge(m, src)
}
func (m *RegisterRequest) XXX_Size() int {
	return xxx_messageInfo_RegisterRequest.Size(m)
}
func (m *RegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterRequest proto.InternalMessageInfo

func (m *RegisterRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *RegisterRequest) GetPosition() *Position {
	if m != nil {
		return m.Position
	}
	return nil
}

func (m *RegisterRequest) GetProperties() []*DetectableProperty {
	if m != nil {
		return m.Properties
	}
	return nil
}

type RegisterResponse struct {
	Effect               string   `protobuf:"bytes,1,opt,name=effect,proto3" json:"effect,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RegisterResponse) Reset()         { *m = RegisterResponse{} }
func (m *RegisterResponse) String() string { return proto.CompactTextString(m) }
func (*RegisterResponse) ProtoMessage()    {}
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{5}
}

func (m *RegisterResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RegisterResponse.Unmarshal(m, b)
}
func (m *RegisterResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RegisterResponse.Marshal(b, m, deterministic)
}
func (m *RegisterResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RegisterResponse.Merge(m, src)
}
func (m *RegisterResponse) XXX_Size() int {
	return xxx_messageInfo_RegisterResponse.Size(m)
}
func (m *RegisterResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RegisterResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RegisterResponse proto.InternalMessageInfo

func (m *RegisterResponse) GetEffect() string {
	if m != nil {
		return m.Effect
	}
	return ""
}

type DetectRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DetectRequest) Reset()         { *m = DetectRequest{} }
func (m *DetectRequest) String() string { return proto.CompactTextString(m) }
func (*DetectRequest) ProtoMessage()    {}
func (*DetectRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{6}
}

func (m *DetectRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DetectRequest.Unmarshal(m, b)
}
func (m *DetectRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DetectRequest.Marshal(b, m, deterministic)
}
func (m *DetectRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DetectRequest.Merge(m, src)
}
func (m *DetectRequest) XXX_Size() int {
	return xxx_messageInfo_DetectRequest.Size(m)
}
func (m *DetectRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DetectRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DetectRequest proto.InternalMessageInfo

type DetectResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DetectResponse) Reset()         { *m = DetectResponse{} }
func (m *DetectResponse) String() string { return proto.CompactTextString(m) }
func (*DetectResponse) ProtoMessage()    {}
func (*DetectResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{7}
}

func (m *DetectResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DetectResponse.Unmarshal(m, b)
}
func (m *DetectResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DetectResponse.Marshal(b, m, deterministic)
}
func (m *DetectResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DetectResponse.Merge(m, src)
}
func (m *DetectResponse) XXX_Size() int {
	return xxx_messageInfo_DetectResponse.Size(m)
}
func (m *DetectResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DetectResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DetectResponse proto.InternalMessageInfo

type GetRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{8}
}

func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (m *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(m, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

type GetResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}
func (*GetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{9}
}

func (m *GetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetResponse.Unmarshal(m, b)
}
func (m *GetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetResponse.Marshal(b, m, deterministic)
}
func (m *GetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetResponse.Merge(m, src)
}
func (m *GetResponse) XXX_Size() int {
	return xxx_messageInfo_GetResponse.Size(m)
}
func (m *GetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetResponse proto.InternalMessageInfo

type InitializeRequest struct {
	Vessel               *Vessel  `protobuf:"bytes,1,opt,name=vessel,proto3" json:"vessel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitializeRequest) Reset()         { *m = InitializeRequest{} }
func (m *InitializeRequest) String() string { return proto.CompactTextString(m) }
func (*InitializeRequest) ProtoMessage()    {}
func (*InitializeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{10}
}

func (m *InitializeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitializeRequest.Unmarshal(m, b)
}
func (m *InitializeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitializeRequest.Marshal(b, m, deterministic)
}
func (m *InitializeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitializeRequest.Merge(m, src)
}
func (m *InitializeRequest) XXX_Size() int {
	return xxx_messageInfo_InitializeRequest.Size(m)
}
func (m *InitializeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InitializeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InitializeRequest proto.InternalMessageInfo

func (m *InitializeRequest) GetVessel() *Vessel {
	if m != nil {
		return m.Vessel
	}
	return nil
}

type InitializeResponse struct {
	Vessel               *Vessel  `protobuf:"bytes,1,opt,name=vessel,proto3" json:"vessel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InitializeResponse) Reset()         { *m = InitializeResponse{} }
func (m *InitializeResponse) String() string { return proto.CompactTextString(m) }
func (*InitializeResponse) ProtoMessage()    {}
func (*InitializeResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a139271552183648, []int{11}
}

func (m *InitializeResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InitializeResponse.Unmarshal(m, b)
}
func (m *InitializeResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InitializeResponse.Marshal(b, m, deterministic)
}
func (m *InitializeResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InitializeResponse.Merge(m, src)
}
func (m *InitializeResponse) XXX_Size() int {
	return xxx_messageInfo_InitializeResponse.Size(m)
}
func (m *InitializeResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_InitializeResponse.DiscardUnknown(m)
}

var xxx_messageInfo_InitializeResponse proto.InternalMessageInfo

func (m *InitializeResponse) GetVessel() *Vessel {
	if m != nil {
		return m.Vessel
	}
	return nil
}

func init() {
	proto.RegisterType((*IdentificationProperty)(nil), "deltav.model.IdentificationProperty")
	proto.RegisterType((*RadiationProperty)(nil), "deltav.model.RadiationProperty")
	proto.RegisterType((*SensorProperty)(nil), "deltav.model.SensorProperty")
	proto.RegisterType((*DetectableProperty)(nil), "deltav.model.DetectableProperty")
	proto.RegisterType((*RegisterRequest)(nil), "deltav.model.RegisterRequest")
	proto.RegisterType((*RegisterResponse)(nil), "deltav.model.RegisterResponse")
	proto.RegisterType((*DetectRequest)(nil), "deltav.model.DetectRequest")
	proto.RegisterType((*DetectResponse)(nil), "deltav.model.DetectResponse")
	proto.RegisterType((*GetRequest)(nil), "deltav.model.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "deltav.model.GetResponse")
	proto.RegisterType((*InitializeRequest)(nil), "deltav.model.InitializeRequest")
	proto.RegisterType((*InitializeResponse)(nil), "deltav.model.InitializeResponse")
}

func init() { proto.RegisterFile("worldmodel.proto", fileDescriptor_a139271552183648) }

var fileDescriptor_a139271552183648 = []byte{
	// 494 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x5d, 0x6f, 0xd3, 0x30,
	0x14, 0x5d, 0x53, 0x54, 0xb5, 0x37, 0x5d, 0xd7, 0x59, 0xa8, 0x0a, 0xd9, 0x60, 0x51, 0xc4, 0x43,
	0x85, 0x50, 0x91, 0x82, 0xc4, 0x13, 0x12, 0x30, 0x0d, 0x8d, 0x0a, 0x21, 0x0d, 0x83, 0xe0, 0x81,
	0x07, 0x94, 0xd5, 0xb7, 0xc8, 0x52, 0x88, 0x83, 0x6d, 0x86, 0xca, 0x1f, 0xe1, 0x77, 0xf0, 0xc4,
	0xdf, 0x43, 0x71, 0xec, 0x7c, 0x6c, 0xed, 0xc3, 0x1e, 0x7d, 0xcf, 0xb9, 0xc7, 0xe7, 0xde, 0xe3,
	0x04, 0xa6, 0xbf, 0x84, 0xcc, 0xd8, 0x77, 0xc1, 0x30, 0x5b, 0x14, 0x52, 0x68, 0x41, 0xc6, 0x0c,
	0x33, 0x9d, 0x5e, 0x2d, 0x4c, 0x2d, 0x9c, 0x14, 0x42, 0x71, 0xcd, 0x45, 0x5e, 0xa1, 0xe1, 0xf8,
	0x0a, 0x95, 0x72, 0xdc, 0xf8, 0x0b, 0xcc, 0x96, 0x0c, 0x73, 0xcd, 0xd7, 0x7c, 0x95, 0x96, 0xac,
	0x0b, 0x29, 0x0a, 0x94, 0x7a, 0x43, 0x22, 0xf0, 0xb5, 0x4c, 0x73, 0x55, 0x88, 0x9c, 0xa1, 0x0c,
	0x7a, 0x51, 0x6f, 0x3e, 0xa2, 0xed, 0x12, 0x39, 0x01, 0x9f, 0xe7, 0x1a, 0x65, 0x9e, 0x66, 0x5f,
	0x39, 0x0b, 0x3c, 0xc3, 0x00, 0x57, 0x5a, 0xb2, 0xf8, 0x0c, 0x0e, 0x69, 0xca, 0x78, 0x57, 0xf7,
	0x09, 0xdc, 0xd1, 0x9b, 0x02, 0x8d, 0xe0, 0x24, 0x39, 0x5a, 0xb4, 0xcd, 0x2e, 0x6a, 0xfa, 0xc7,
	0x4d, 0x81, 0xd4, 0x10, 0xe3, 0x29, 0x4c, 0x3e, 0x60, 0xae, 0x84, 0x74, 0x12, 0xf1, 0xbf, 0x1e,
	0x90, 0x33, 0xd4, 0xb8, 0xd2, 0xe9, 0x65, 0x86, 0xb5, 0xf2, 0x33, 0xf0, 0x38, 0x33, 0xba, 0x7e,
	0xf2, 0xb0, 0xab, 0xbb, 0x7d, 0xc6, 0x37, 0x7b, 0xd4, 0xe3, 0x8c, 0xbc, 0x80, 0x91, 0x74, 0xf7,
	0x9a, 0x29, 0xfc, 0xe4, 0x64, 0x87, 0xad, 0x56, 0x67, 0xd3, 0x43, 0x8e, 0x61, 0x54, 0x4e, 0x9d,
	0x2b, 0xae, 0x37, 0x01, 0x8b, 0x7a, 0x73, 0x8f, 0x36, 0x85, 0x53, 0x80, 0x61, 0xe1, 0x9c, 0xff,
	0xe9, 0xc1, 0x01, 0xc5, 0x6f, 0x5c, 0x69, 0x94, 0x14, 0x7f, 0xfc, 0x44, 0xa5, 0xc9, 0xa4, 0xb6,
	0x3d, 0x32, 0x76, 0x12, 0x18, 0xba, 0xc8, 0xac, 0x9b, 0x59, 0xd7, 0xcd, 0x85, 0x45, 0x69, 0xcd,
	0x23, 0x2f, 0x01, 0xec, 0x1d, 0x1c, 0x55, 0xd0, 0x8f, 0xfa, 0x73, 0x3f, 0x89, 0xba, 0x5d, 0x37,
	0x17, 0x46, 0x5b, 0x3d, 0xf1, 0x23, 0x98, 0x36, 0xc6, 0xca, 0x84, 0x15, 0x92, 0x19, 0x0c, 0x70,
	0xbd, 0xc6, 0x95, 0xb6, 0xee, 0xec, 0x29, 0x3e, 0x80, 0xfd, 0x4a, 0xcd, 0x8e, 0x50, 0x46, 0xe4,
	0x0a, 0x55, 0x6b, 0x3c, 0x06, 0x38, 0xc7, 0x1a, 0xdf, 0x07, 0xdf, 0x9c, 0x2c, 0xf8, 0x0a, 0x0e,
	0x97, 0x39, 0xd7, 0x3c, 0xcd, 0xf8, 0x6f, 0x74, 0x6b, 0x78, 0x0c, 0x83, 0xea, 0x65, 0xda, 0x04,
	0xef, 0x76, 0xed, 0x7f, 0x32, 0x18, 0xb5, 0x9c, 0xf8, 0x14, 0x48, 0x5b, 0xc2, 0x1a, 0xbe, 0x95,
	0x46, 0xf2, 0xd7, 0x03, 0xf8, 0x5c, 0x7e, 0x3c, 0xef, 0x4a, 0x94, 0xbc, 0x85, 0xa1, 0xdb, 0x00,
	0xb9, 0x7f, 0x2d, 0xff, 0x6e, 0x64, 0xe1, 0x83, 0x5d, 0xb0, 0x1d, 0x70, 0x8f, 0xbc, 0x86, 0x41,
	0xb5, 0x11, 0x72, 0xb4, 0x2d, 0x06, 0x27, 0x74, 0xbc, 0x1d, 0xac, 0x65, 0x9e, 0x43, 0xff, 0x1c,
	0x35, 0x09, 0xba, 0xb4, 0x66, 0xb3, 0xe1, 0xbd, 0x2d, 0x48, 0xdd, 0xfd, 0x1e, 0xa0, 0x59, 0x12,
	0xb9, 0xf6, 0xa6, 0x6f, 0x24, 0x10, 0x46, 0xbb, 0x09, 0x4e, 0xf2, 0x72, 0x60, 0x7e, 0x1b, 0x4f,
	0xff, 0x07, 0x00, 0x00, 0xff, 0xff, 0x43, 0x67, 0x09, 0xb6, 0x76, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WorldModelClient is the client API for WorldModel service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WorldModelClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Detect(ctx context.Context, in *DetectRequest, opts ...grpc.CallOption) (*DetectResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error)
	Initialize(ctx context.Context, in *InitializeRequest, opts ...grpc.CallOption) (*InitializeResponse, error)
}

type worldModelClient struct {
	cc *grpc.ClientConn
}

func NewWorldModelClient(cc *grpc.ClientConn) WorldModelClient {
	return &worldModelClient{cc}
}

func (c *worldModelClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/deltav.model.WorldModel/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *worldModelClient) Detect(ctx context.Context, in *DetectRequest, opts ...grpc.CallOption) (*DetectResponse, error) {
	out := new(DetectResponse)
	err := c.cc.Invoke(ctx, "/deltav.model.WorldModel/Detect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *worldModelClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	out := new(GetResponse)
	err := c.cc.Invoke(ctx, "/deltav.model.WorldModel/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *worldModelClient) Initialize(ctx context.Context, in *InitializeRequest, opts ...grpc.CallOption) (*InitializeResponse, error) {
	out := new(InitializeResponse)
	err := c.cc.Invoke(ctx, "/deltav.model.WorldModel/Initialize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WorldModelServer is the server API for WorldModel service.
type WorldModelServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Detect(context.Context, *DetectRequest) (*DetectResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	Initialize(context.Context, *InitializeRequest) (*InitializeResponse, error)
}

func RegisterWorldModelServer(s *grpc.Server, srv WorldModelServer) {
	s.RegisterService(&_WorldModel_serviceDesc, srv)
}

func _WorldModel_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorldModelServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/deltav.model.WorldModel/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorldModelServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorldModel_Detect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorldModelServer).Detect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/deltav.model.WorldModel/Detect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorldModelServer).Detect(ctx, req.(*DetectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorldModel_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorldModelServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/deltav.model.WorldModel/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorldModelServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _WorldModel_Initialize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitializeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WorldModelServer).Initialize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/deltav.model.WorldModel/Initialize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WorldModelServer).Initialize(ctx, req.(*InitializeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _WorldModel_serviceDesc = grpc.ServiceDesc{
	ServiceName: "deltav.model.WorldModel",
	HandlerType: (*WorldModelServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _WorldModel_Register_Handler,
		},
		{
			MethodName: "Detect",
			Handler:    _WorldModel_Detect_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _WorldModel_Get_Handler,
		},
		{
			MethodName: "Initialize",
			Handler:    _WorldModel_Initialize_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "worldmodel.proto",
}