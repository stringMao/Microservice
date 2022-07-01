// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: commod.proto

package base

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

//
type Txt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txt string `protobuf:"bytes,1,opt,name=txt,proto3" json:"txt,omitempty"`
}

func (x *Txt) Reset() {
	*x = Txt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commod_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Txt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Txt) ProtoMessage() {}

func (x *Txt) ProtoReflect() protoreflect.Message {
	mi := &file_commod_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Txt.ProtoReflect.Descriptor instead.
func (*Txt) Descriptor() ([]byte, []int) {
	return file_commod_proto_rawDescGZIP(), []int{0}
}

func (x *Txt) GetTxt() string {
	if x != nil {
		return x.Txt
	}
	return ""
}

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commod_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_commod_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_commod_proto_rawDescGZIP(), []int{1}
}

func (x *Status) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type ReplyResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Txt  string `protobuf:"bytes,2,opt,name=txt,proto3" json:"txt,omitempty"`
}

func (x *ReplyResult) Reset() {
	*x = ReplyResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commod_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplyResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplyResult) ProtoMessage() {}

func (x *ReplyResult) ProtoReflect() protoreflect.Message {
	mi := &file_commod_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplyResult.ProtoReflect.Descriptor instead.
func (*ReplyResult) Descriptor() ([]byte, []int) {
	return file_commod_proto_rawDescGZIP(), []int{2}
}

func (x *ReplyResult) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ReplyResult) GetTxt() string {
	if x != nil {
		return x.Txt
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageId uint32 `protobuf:"varint,1,opt,name=messageId,proto3" json:"messageId,omitempty"`
	Body      []byte `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commod_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_commod_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_commod_proto_rawDescGZIP(), []int{3}
}

func (x *Message) GetMessageId() uint32 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

func (x *Message) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

type Forward struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId   uint64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	ServerId uint64 `protobuf:"varint,2,opt,name=ServerId,proto3" json:"ServerId,omitempty"`
	Body     []byte `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Forward) Reset() {
	*x = Forward{}
	if protoimpl.UnsafeEnabled {
		mi := &file_commod_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Forward) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Forward) ProtoMessage() {}

func (x *Forward) ProtoReflect() protoreflect.Message {
	mi := &file_commod_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Forward.ProtoReflect.Descriptor instead.
func (*Forward) Descriptor() ([]byte, []int) {
	return file_commod_proto_rawDescGZIP(), []int{4}
}

func (x *Forward) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Forward) GetServerId() uint64 {
	if x != nil {
		return x.ServerId
	}
	return 0
}

func (x *Forward) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

var File_commod_proto protoreflect.FileDescriptor

var file_commod_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04,
	0x62, 0x61, 0x73, 0x65, 0x22, 0x17, 0x0a, 0x03, 0x54, 0x78, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x74,
	0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x78, 0x74, 0x22, 0x1c, 0x0a,
	0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x33, 0x0a, 0x0b, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x74, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x78, 0x74,
	0x22, 0x3b, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x64,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x51, 0x0a,
	0x07, 0x46, 0x6f, 0x72, 0x77, 0x61, 0x72, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x62, 0x6f, 0x64, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x42, 0x07, 0x5a, 0x05, 0x62, 0x61, 0x73, 0x65, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_commod_proto_rawDescOnce sync.Once
	file_commod_proto_rawDescData = file_commod_proto_rawDesc
)

func file_commod_proto_rawDescGZIP() []byte {
	file_commod_proto_rawDescOnce.Do(func() {
		file_commod_proto_rawDescData = protoimpl.X.CompressGZIP(file_commod_proto_rawDescData)
	})
	return file_commod_proto_rawDescData
}

var file_commod_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_commod_proto_goTypes = []interface{}{
	(*Txt)(nil),         // 0: base.Txt
	(*Status)(nil),      // 1: base.Status
	(*ReplyResult)(nil), // 2: base.ReplyResult
	(*Message)(nil),     // 3: base.Message
	(*Forward)(nil),     // 4: base.Forward
}
var file_commod_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_commod_proto_init() }
func file_commod_proto_init() {
	if File_commod_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_commod_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Txt); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commod_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commod_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplyResult); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commod_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_commod_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Forward); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_commod_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_commod_proto_goTypes,
		DependencyIndexes: file_commod_proto_depIdxs,
		MessageInfos:      file_commod_proto_msgTypes,
	}.Build()
	File_commod_proto = out.File
	file_commod_proto_rawDesc = nil
	file_commod_proto_goTypes = nil
	file_commod_proto_depIdxs = nil
}
