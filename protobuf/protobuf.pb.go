// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf.proto

package protobuf

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Kmessage struct {
	Label                string   `protobuf:"bytes,1,opt,name=Label,proto3" json:"Label,omitempty"`
	SenderId             string   `protobuf:"bytes,2,opt,name=SenderId,proto3" json:"SenderId,omitempty"`
	SenderAddr           string   `protobuf:"bytes,3,opt,name=SenderAddr,proto3" json:"SenderAddr,omitempty"`
	LookupId             string   `protobuf:"bytes,4,opt,name=LookupId,proto3" json:"LookupId,omitempty"`
	LookupResp           string   `protobuf:"bytes,5,opt,name=LookupResp,proto3" json:"LookupResp,omitempty"`
	Key                  string   `protobuf:"bytes,6,opt,name=Key,proto3" json:"Key,omitempty"`
	Data                 string   `protobuf:"bytes,7,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Kmessage) Reset()         { *m = Kmessage{} }
func (m *Kmessage) String() string { return proto.CompactTextString(m) }
func (*Kmessage) ProtoMessage()    {}
func (*Kmessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_c77a803fcbc0c059, []int{0}
}

func (m *Kmessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Kmessage.Unmarshal(m, b)
}
func (m *Kmessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Kmessage.Marshal(b, m, deterministic)
}
func (m *Kmessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Kmessage.Merge(m, src)
}
func (m *Kmessage) XXX_Size() int {
	return xxx_messageInfo_Kmessage.Size(m)
}
func (m *Kmessage) XXX_DiscardUnknown() {
	xxx_messageInfo_Kmessage.DiscardUnknown(m)
}

var xxx_messageInfo_Kmessage proto.InternalMessageInfo

func (m *Kmessage) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *Kmessage) GetSenderId() string {
	if m != nil {
		return m.SenderId
	}
	return ""
}

func (m *Kmessage) GetSenderAddr() string {
	if m != nil {
		return m.SenderAddr
	}
	return ""
}

func (m *Kmessage) GetLookupId() string {
	if m != nil {
		return m.LookupId
	}
	return ""
}

func (m *Kmessage) GetLookupResp() string {
	if m != nil {
		return m.LookupResp
	}
	return ""
}

func (m *Kmessage) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Kmessage) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*Kmessage)(nil), "protobuf.Kmessage")
}

func init() { proto.RegisterFile("protobuf.proto", fileDescriptor_c77a803fcbc0c059) }

var fileDescriptor_c77a803fcbc0c059 = []byte{
	// 166 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x4f, 0x2a, 0x4d, 0xd3, 0x03, 0x33, 0x84, 0x38, 0x60, 0x7c, 0xa5, 0x7d, 0x8c, 0x5c, 0x1c,
	0xde, 0xb9, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x42, 0x22, 0x5c, 0xac, 0x3e, 0x89, 0x49, 0xa9,
	0x39, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x10, 0x8e, 0x90, 0x14, 0x17, 0x47, 0x70, 0x6a,
	0x5e, 0x4a, 0x6a, 0x91, 0x67, 0x8a, 0x04, 0x13, 0x58, 0x02, 0xce, 0x17, 0x92, 0xe3, 0xe2, 0x82,
	0xb0, 0x1d, 0x53, 0x52, 0x8a, 0x24, 0x98, 0xc1, 0xb2, 0x48, 0x22, 0x20, 0xbd, 0x3e, 0xf9, 0xf9,
	0xd9, 0xa5, 0x05, 0x9e, 0x29, 0x12, 0x2c, 0x10, 0xbd, 0x30, 0x3e, 0x48, 0x2f, 0x84, 0x1d, 0x94,
	0x5a, 0x5c, 0x20, 0xc1, 0x0a, 0xd1, 0x8b, 0x10, 0x11, 0x12, 0xe0, 0x62, 0xf6, 0x4e, 0xad, 0x94,
	0x60, 0x03, 0x4b, 0x80, 0x98, 0x42, 0x42, 0x5c, 0x2c, 0x2e, 0x89, 0x25, 0x89, 0x12, 0xec, 0x60,
	0x21, 0x30, 0x3b, 0x89, 0x0d, 0xec, 0x15, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x7c, 0x43,
	0x72, 0xbd, 0xe3, 0x00, 0x00, 0x00,
}
