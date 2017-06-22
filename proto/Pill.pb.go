// Code generated by protoc-gen-go.
// source: Pill.proto
// DO NOT EDIT!

/*
Package Proto is a generated protocol buffer package.

It is generated from these files:
	Pill.proto

It has these top-level messages:
	ChatHeader
	ChatData
	ChatList
	BlockData
	WorkerHandShark
	MessageData
	HandShake
*/
package Proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type ChatHeader struct {
	// optional int32 RoomId   = 1;
	// optional int32 Type   = 2;
	// optional int32 Status   = 3;
	// optional int32 Cmd   = 4;
	// optional int32 Sid   = 5;
	// 正常消息为1，返回消息非1
	Code             *int32  `protobuf:"varint,1,opt,name=Code" json:"Code,omitempty"`
	Msg              *string `protobuf:"bytes,2,opt,name=Msg" json:"Msg,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *ChatHeader) Reset()                    { *m = ChatHeader{} }
func (m *ChatHeader) String() string            { return proto.CompactTextString(m) }
func (*ChatHeader) ProtoMessage()               {}
func (*ChatHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ChatHeader) GetCode() int32 {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return 0
}

func (m *ChatHeader) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

type ChatData struct {
	Header           *ChatHeader `protobuf:"bytes,1,opt,name=Header" json:"Header,omitempty"`
	Msg              *string     `protobuf:"bytes,3,opt,name=Msg" json:"Msg,omitempty"`
	Msgjson          *string     `protobuf:"bytes,4,opt,name=Msgjson" json:"Msgjson,omitempty"`
	Uid              *int32      `protobuf:"varint,5,opt,name=Uid" json:"Uid,omitempty"`
	Uname            *string     `protobuf:"bytes,6,opt,name=Uname" json:"Uname,omitempty"`
	Upic             *string     `protobuf:"bytes,7,opt,name=Upic" json:"Upic,omitempty"`
	Utitle           *string     `protobuf:"bytes,8,opt,name=Utitle" json:"Utitle,omitempty"`
	Img              *string     `protobuf:"bytes,9,opt,name=Img" json:"Img,omitempty"`
	Timestamp        *int32      `protobuf:"varint,10,opt,name=Timestamp" json:"Timestamp,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *ChatData) Reset()                    { *m = ChatData{} }
func (m *ChatData) String() string            { return proto.CompactTextString(m) }
func (*ChatData) ProtoMessage()               {}
func (*ChatData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ChatData) GetHeader() *ChatHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *ChatData) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

func (m *ChatData) GetMsgjson() string {
	if m != nil && m.Msgjson != nil {
		return *m.Msgjson
	}
	return ""
}

func (m *ChatData) GetUid() int32 {
	if m != nil && m.Uid != nil {
		return *m.Uid
	}
	return 0
}

func (m *ChatData) GetUname() string {
	if m != nil && m.Uname != nil {
		return *m.Uname
	}
	return ""
}

func (m *ChatData) GetUpic() string {
	if m != nil && m.Upic != nil {
		return *m.Upic
	}
	return ""
}

func (m *ChatData) GetUtitle() string {
	if m != nil && m.Utitle != nil {
		return *m.Utitle
	}
	return ""
}

func (m *ChatData) GetImg() string {
	if m != nil && m.Img != nil {
		return *m.Img
	}
	return ""
}

func (m *ChatData) GetTimestamp() int32 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

type ChatList struct {
	ChatItem         []*ChatData `protobuf:"bytes,1,rep,name=ChatItem" json:"ChatItem,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *ChatList) Reset()                    { *m = ChatList{} }
func (m *ChatList) String() string            { return proto.CompactTextString(m) }
func (*ChatList) ProtoMessage()               {}
func (*ChatList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ChatList) GetChatItem() []*ChatData {
	if m != nil {
		return m.ChatItem
	}
	return nil
}

type BlockData struct {
	Header           *ChatHeader `protobuf:"bytes,1,opt,name=Header" json:"Header,omitempty"`
	BlockUid         *int32      `protobuf:"varint,2,opt,name=BlockUid" json:"BlockUid,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *BlockData) Reset()                    { *m = BlockData{} }
func (m *BlockData) String() string            { return proto.CompactTextString(m) }
func (*BlockData) ProtoMessage()               {}
func (*BlockData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *BlockData) GetHeader() *ChatHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *BlockData) GetBlockUid() int32 {
	if m != nil && m.BlockUid != nil {
		return *m.BlockUid
	}
	return 0
}

type WorkerHandShark struct {
	IP               *string `protobuf:"bytes,1,opt,name=IP" json:"IP,omitempty"`
	Name             *string `protobuf:"bytes,2,opt,name=Name" json:"Name,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *WorkerHandShark) Reset()                    { *m = WorkerHandShark{} }
func (m *WorkerHandShark) String() string            { return proto.CompactTextString(m) }
func (*WorkerHandShark) ProtoMessage()               {}
func (*WorkerHandShark) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *WorkerHandShark) GetIP() string {
	if m != nil && m.IP != nil {
		return *m.IP
	}
	return ""
}

func (m *WorkerHandShark) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

type MessageData struct {
	Header           *ChatHeader `protobuf:"bytes,1,opt,name=Header" json:"Header,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *MessageData) Reset()                    { *m = MessageData{} }
func (m *MessageData) String() string            { return proto.CompactTextString(m) }
func (*MessageData) ProtoMessage()               {}
func (*MessageData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *MessageData) GetHeader() *ChatHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type HandShake struct {
	Header           *ChatHeader `protobuf:"bytes,1,opt,name=Header" json:"Header,omitempty"`
	Timestamp        *int32      `protobuf:"varint,2,opt,name=Timestamp" json:"Timestamp,omitempty"`
	Token            *uint32     `protobuf:"varint,3,opt,name=Token" json:"Token,omitempty"`
	RoomId           *int32      `protobuf:"varint,4,opt,name=RoomId" json:"RoomId,omitempty"`
	ConnectId        *int32      `protobuf:"varint,5,opt,name=ConnectId" json:"ConnectId,omitempty"`
	Uid              *int32      `protobuf:"varint,6,opt,name=Uid" json:"Uid,omitempty"`
	Platform         *int32      `protobuf:"varint,7,opt,name=Platform" json:"Platform,omitempty"`
	Uname            *string     `protobuf:"bytes,8,opt,name=Uname" json:"Uname,omitempty"`
	XXX_unrecognized []byte      `json:"-"`
}

func (m *HandShake) Reset()                    { *m = HandShake{} }
func (m *HandShake) String() string            { return proto.CompactTextString(m) }
func (*HandShake) ProtoMessage()               {}
func (*HandShake) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *HandShake) GetHeader() *ChatHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *HandShake) GetTimestamp() int32 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *HandShake) GetToken() uint32 {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return 0
}

func (m *HandShake) GetRoomId() int32 {
	if m != nil && m.RoomId != nil {
		return *m.RoomId
	}
	return 0
}

func (m *HandShake) GetConnectId() int32 {
	if m != nil && m.ConnectId != nil {
		return *m.ConnectId
	}
	return 0
}

func (m *HandShake) GetUid() int32 {
	if m != nil && m.Uid != nil {
		return *m.Uid
	}
	return 0
}

func (m *HandShake) GetPlatform() int32 {
	if m != nil && m.Platform != nil {
		return *m.Platform
	}
	return 0
}

func (m *HandShake) GetUname() string {
	if m != nil && m.Uname != nil {
		return *m.Uname
	}
	return ""
}

func init() {
	proto.RegisterType((*ChatHeader)(nil), "Proto.ChatHeader")
	proto.RegisterType((*ChatData)(nil), "Proto.ChatData")
	proto.RegisterType((*ChatList)(nil), "Proto.ChatList")
	proto.RegisterType((*BlockData)(nil), "Proto.BlockData")
	proto.RegisterType((*WorkerHandShark)(nil), "Proto.WorkerHandShark")
	proto.RegisterType((*MessageData)(nil), "Proto.MessageData")
	proto.RegisterType((*HandShake)(nil), "Proto.HandShake")
}

var fileDescriptor0 = []byte{
	// 385 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x94, 0x52, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x25, 0x69, 0x93, 0x26, 0x53, 0xb4, 0xba, 0x88, 0x2c, 0xe2, 0x41, 0x72, 0x52, 0x84, 0x1e,
	0x0a, 0xa2, 0x67, 0xeb, 0xa1, 0x01, 0x2b, 0x61, 0x6d, 0xf0, 0xbc, 0x34, 0x6b, 0x1b, 0xf3, 0xb1,
	0x25, 0xd9, 0x9f, 0xac, 0xff, 0xc3, 0xdd, 0x49, 0xd2, 0x54, 0x6f, 0x3d, 0xed, 0x7b, 0xc3, 0xce,
	0xdb, 0x99, 0xb7, 0x0f, 0x20, 0x4a, 0xf3, 0x7c, 0xba, 0xab, 0xa4, 0x92, 0xc4, 0x89, 0xcc, 0x11,
	0xcc, 0x00, 0xe6, 0x5b, 0xae, 0x16, 0x82, 0x27, 0xa2, 0x22, 0x04, 0x86, 0x73, 0x99, 0x08, 0x6a,
	0xdd, 0x58, 0xb7, 0x0e, 0x43, 0x4c, 0xce, 0x60, 0xb0, 0xac, 0x37, 0xd4, 0xd6, 0x25, 0x9f, 0x19,
	0x18, 0xfc, 0x58, 0xe0, 0x99, 0xa6, 0x17, 0xae, 0x38, 0xb9, 0x03, 0xb7, 0x69, 0xc6, 0xa6, 0xf1,
	0xec, 0x7c, 0x8a, 0xc2, 0xd3, 0x5e, 0x95, 0xb5, 0x17, 0x3a, 0xa5, 0xc1, 0x5e, 0x89, 0x50, 0x18,
	0xe9, 0xe3, 0xab, 0x96, 0x25, 0x1d, 0x62, 0xb5, 0xa3, 0xe6, 0x6e, 0x9c, 0x26, 0xd4, 0xc1, 0x41,
	0x0c, 0x24, 0x17, 0xe0, 0xc4, 0x25, 0x2f, 0x04, 0x75, 0xf1, 0x66, 0x43, 0xcc, 0xc4, 0xf1, 0x2e,
	0x5d, 0xd3, 0x11, 0x16, 0x11, 0x93, 0x4b, 0x70, 0x63, 0x95, 0xaa, 0x5c, 0x50, 0x0f, 0xab, 0x2d,
	0x33, 0x9a, 0x61, 0xb1, 0xa1, 0x7e, 0xf3, 0xbe, 0x86, 0xe4, 0x1a, 0xfc, 0x55, 0x5a, 0x88, 0x5a,
	0xf1, 0x62, 0x47, 0x01, 0xdf, 0xea, 0x0b, 0xc1, 0x63, 0xb3, 0xe6, 0x6b, 0x5a, 0x2b, 0x72, 0xdf,
	0xe0, 0x50, 0x89, 0x42, 0x2f, 0x3a, 0xd0, 0x8b, 0x4e, 0x0e, 0x16, 0x35, 0x4e, 0xb0, 0xfd, 0x85,
	0x80, 0x81, 0xff, 0x9c, 0xcb, 0x75, 0x76, 0xac, 0x41, 0x57, 0xe0, 0x61, 0x9f, 0xd9, 0xdc, 0xc6,
	0x69, 0xf6, 0x3c, 0x78, 0x80, 0xc9, 0x87, 0xac, 0x32, 0x51, 0x2d, 0x78, 0x99, 0xbc, 0x6f, 0x79,
	0x95, 0x91, 0x53, 0xb0, 0xc3, 0x08, 0x55, 0x7d, 0xa6, 0x91, 0xf1, 0xe2, 0xcd, 0x18, 0xd4, 0x7c,
	0x15, 0xe2, 0xe0, 0x09, 0xc6, 0x4b, 0x51, 0xd7, 0x7c, 0x23, 0x8e, 0x1c, 0x26, 0xf8, 0xb6, 0xc0,
	0x6f, 0xdf, 0xca, 0xc4, 0x31, 0x5b, 0xfc, 0x31, 0xd5, 0xfe, 0x67, 0xaa, 0xf9, 0xc6, 0x95, 0xcc,
	0x44, 0x89, 0x31, 0x38, 0x61, 0x0d, 0x31, 0x5f, 0xc6, 0xa4, 0x2c, 0xc2, 0x04, 0x73, 0xe0, 0xb0,
	0x96, 0x19, 0xad, 0xb9, 0x2c, 0x4b, 0xb1, 0x56, 0x61, 0x17, 0x86, 0xbe, 0xd0, 0x85, 0xc4, 0xed,
	0x43, 0xa2, 0x1d, 0x8c, 0x72, 0xae, 0x3e, 0x65, 0x55, 0x60, 0x24, 0xb4, 0x83, 0x1d, 0xef, 0x03,
	0xe4, 0x1d, 0x04, 0xe8, 0x37, 0x00, 0x00, 0xff, 0xff, 0x4f, 0xd3, 0xde, 0x5f, 0x14, 0x03, 0x00,
	0x00,
}
