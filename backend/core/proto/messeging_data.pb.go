// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: messeging_data.proto

package proto

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

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header map[string]string `protobuf:"bytes,1,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body   *Body             `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messeging_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_messeging_data_proto_msgTypes[0]
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
	return file_messeging_data_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetHeader() map[string]string {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *Message) GetBody() *Body {
	if x != nil {
		return x.Body
	}
	return nil
}

type Body struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Payload:
	//
	//	*Body_Trigger
	//	*Body_Chunking
	Payload isBody_Payload `protobuf_oneof:"payload"`
}

func (x *Body) Reset() {
	*x = Body{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messeging_data_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Body) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Body) ProtoMessage() {}

func (x *Body) ProtoReflect() protoreflect.Message {
	mi := &file_messeging_data_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Body.ProtoReflect.Descriptor instead.
func (*Body) Descriptor() ([]byte, []int) {
	return file_messeging_data_proto_rawDescGZIP(), []int{1}
}

func (m *Body) GetPayload() isBody_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (x *Body) GetTrigger() *ConnectorRequest {
	if x, ok := x.GetPayload().(*Body_Trigger); ok {
		return x.Trigger
	}
	return nil
}

func (x *Body) GetChunking() *ChunkingData {
	if x, ok := x.GetPayload().(*Body_Chunking); ok {
		return x.Chunking
	}
	return nil
}

type isBody_Payload interface {
	isBody_Payload()
}

type Body_Trigger struct {
	Trigger *ConnectorRequest `protobuf:"bytes,1,opt,name=trigger,proto3,oneof"`
}

type Body_Chunking struct {
	Chunking *ChunkingData `protobuf:"bytes,2,opt,name=chunking,proto3,oneof"`
}

func (*Body_Trigger) isBody_Payload() {}

func (*Body_Chunking) isBody_Payload() {}

var File_messeging_data_proto protoreflect.FileDescriptor

var file_messeging_data_proto_rawDesc = []byte{
	0x0a, 0x14, 0x6d, 0x65, 0x73, 0x73, 0x65, 0x67, 0x69, 0x6e, 0x67, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6d, 0x62, 0x65,
	0x64, 0x64, 0x1a, 0x18, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x13, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xa3, 0x01, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x37, 0x0a,
	0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x24, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64,
	0x64, 0x2e, 0x42, 0x6f, 0x64, 0x79, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x1a, 0x39, 0x0a, 0x0b,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x83, 0x01, 0x0a, 0x04, 0x42, 0x6f, 0x64, 0x79,
	0x12, 0x38, 0x0a, 0x07, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x2e, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48,
	0x00, 0x52, 0x07, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12, 0x36, 0x0a, 0x08, 0x63, 0x68,
	0x75, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x65, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x2e, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x69,
	0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x08, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x69,
	0x6e, 0x67, 0x42, 0x09, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x1a, 0x5a,
	0x18, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_messeging_data_proto_rawDescOnce sync.Once
	file_messeging_data_proto_rawDescData = file_messeging_data_proto_rawDesc
)

func file_messeging_data_proto_rawDescGZIP() []byte {
	file_messeging_data_proto_rawDescOnce.Do(func() {
		file_messeging_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_messeging_data_proto_rawDescData)
	})
	return file_messeging_data_proto_rawDescData
}

var file_messeging_data_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_messeging_data_proto_goTypes = []interface{}{
	(*Message)(nil),          // 0: com.embedd.Message
	(*Body)(nil),             // 1: com.embedd.Body
	nil,                      // 2: com.embedd.Message.HeaderEntry
	(*ConnectorRequest)(nil), // 3: com.embedd.ConnectorRequest
	(*ChunkingData)(nil),     // 4: com.embedd.ChunkingData
}
var file_messeging_data_proto_depIdxs = []int32{
	2, // 0: com.embedd.Message.header:type_name -> com.embedd.Message.HeaderEntry
	1, // 1: com.embedd.Message.body:type_name -> com.embedd.Body
	3, // 2: com.embedd.Body.trigger:type_name -> com.embedd.ConnectorRequest
	4, // 3: com.embedd.Body.chunking:type_name -> com.embedd.ChunkingData
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_messeging_data_proto_init() }
func file_messeging_data_proto_init() {
	if File_messeging_data_proto != nil {
		return
	}
	file_connector_messages_proto_init()
	file_chunking_data_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_messeging_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_messeging_data_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Body); i {
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
	file_messeging_data_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*Body_Trigger)(nil),
		(*Body_Chunking)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messeging_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_messeging_data_proto_goTypes,
		DependencyIndexes: file_messeging_data_proto_depIdxs,
		MessageInfos:      file_messeging_data_proto_msgTypes,
	}.Build()
	File_messeging_data_proto = out.File
	file_messeging_data_proto_rawDesc = nil
	file_messeging_data_proto_goTypes = nil
	file_messeging_data_proto_depIdxs = nil
}