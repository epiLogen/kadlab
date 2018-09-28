// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protobuf.proto

package d7024e

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
	Label                string   `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	SenderId             string   `protobuf:"bytes,2,opt,name=senderId,proto3" json:"senderId,omitempty"`
	SenderAddr           string   `protobuf:"bytes,3,opt,name=senderAddr,proto3" json:"senderAddr,omitempty"`
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

func init() {
	proto.RegisterType((*Kmessage)(nil), "d7024e.Kmessage")
}

func init() { proto.RegisterFile("protobuf.proto", fileDescriptor_c77a803fcbc0c059) }

var fileDescriptor_c77a803fcbc0c059 = []byte{
	// 113 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x4f, 0x2a, 0x4d, 0xd3, 0x03, 0x33, 0x84, 0xd8, 0x52, 0xcc, 0x0d, 0x8c, 0x4c, 0x52, 0x95,
	0x62, 0xb8, 0x38, 0xbc, 0x73, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x85, 0x44, 0xb8, 0x58, 0x73,
	0x12, 0x93, 0x52, 0x73, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x20, 0x1c, 0x21, 0x29, 0x2e,
	0x8e, 0xe2, 0xd4, 0xbc, 0x94, 0xd4, 0x22, 0xcf, 0x14, 0x09, 0x26, 0xb0, 0x04, 0x9c, 0x2f, 0x24,
	0xc7, 0xc5, 0x05, 0x61, 0x3b, 0xa6, 0xa4, 0x14, 0x49, 0x30, 0x83, 0x65, 0x91, 0x44, 0x92, 0xd8,
	0xc0, 0x96, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc8, 0x65, 0x73, 0x9b, 0x7e, 0x00, 0x00,
	0x00,
}
