// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.3
// source: proto/timeline/timeline.proto

package timeline

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	post "grpc-test/protogen/post"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Event int32

const (
	Event_INITIAL_ACCESS Event = 0
	Event_POST_CREATED   Event = 1
	Event_POSTS_DELETED  Event = 2
)

// Enum value maps for Event.
var (
	Event_name = map[int32]string{
		0: "INITIAL_ACCESS",
		1: "POST_CREATED",
		2: "POSTS_DELETED",
	}
	Event_value = map[string]int32{
		"INITIAL_ACCESS": 0,
		"POST_CREATED":   1,
		"POSTS_DELETED":  2,
	}
)

func (x Event) Enum() *Event {
	p := new(Event)
	*p = x
	return p
}

func (x Event) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Event) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_timeline_timeline_proto_enumTypes[0].Descriptor()
}

func (Event) Type() protoreflect.EnumType {
	return &file_proto_timeline_timeline_proto_enumTypes[0]
}

func (x Event) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Event.Descriptor instead.
func (Event) EnumDescriptor() ([]byte, []int) {
	return file_proto_timeline_timeline_proto_rawDescGZIP(), []int{0}
}

type TimelineRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *TimelineRequest) Reset() {
	*x = TimelineRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_timeline_timeline_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimelineRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimelineRequest) ProtoMessage() {}

func (x *TimelineRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_timeline_timeline_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimelineRequest.ProtoReflect.Descriptor instead.
func (*TimelineRequest) Descriptor() ([]byte, []int) {
	return file_proto_timeline_timeline_proto_rawDescGZIP(), []int{0}
}

func (x *TimelineRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type TimelineResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType Event        `protobuf:"varint,1,opt,name=event_type,json=eventType,proto3,enum=Event" json:"event_type,omitempty"`
	Posts     []*post.Post `protobuf:"bytes,2,rep,name=posts,proto3" json:"posts,omitempty"`
}

func (x *TimelineResponse) Reset() {
	*x = TimelineResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_timeline_timeline_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimelineResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimelineResponse) ProtoMessage() {}

func (x *TimelineResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_timeline_timeline_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimelineResponse.ProtoReflect.Descriptor instead.
func (*TimelineResponse) Descriptor() ([]byte, []int) {
	return file_proto_timeline_timeline_proto_rawDescGZIP(), []int{1}
}

func (x *TimelineResponse) GetEventType() Event {
	if x != nil {
		return x.EventType
	}
	return Event_INITIAL_ACCESS
}

func (x *TimelineResponse) GetPosts() []*post.Post {
	if x != nil {
		return x.Posts
	}
	return nil
}

var File_proto_timeline_timeline_proto protoreflect.FileDescriptor

var file_proto_timeline_timeline_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x6f, 0x73, 0x74, 0x2f, 0x70, 0x6f, 0x73, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x21, 0x0a, 0x0f, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69,
	0x6e, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x56, 0x0a, 0x10, 0x54, 0x69, 0x6d,
	0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25, 0x0a,
	0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x06, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x05, 0x70, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x02, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x52, 0x05, 0x70, 0x6f, 0x73, 0x74,
	0x73, 0x2a, 0x40, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x0e, 0x49, 0x4e,
	0x49, 0x54, 0x49, 0x41, 0x4c, 0x5f, 0x41, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x10,
	0x0a, 0x0c, 0x50, 0x4f, 0x53, 0x54, 0x5f, 0x43, 0x52, 0x45, 0x41, 0x54, 0x45, 0x44, 0x10, 0x01,
	0x12, 0x11, 0x0a, 0x0d, 0x50, 0x4f, 0x53, 0x54, 0x53, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45,
	0x44, 0x10, 0x02, 0x32, 0x44, 0x0a, 0x0f, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x31, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x73,
	0x74, 0x73, 0x12, 0x10, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x1d, 0x5a, 0x1b, 0x67, 0x72, 0x70,
	0x63, 0x2d, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_timeline_timeline_proto_rawDescOnce sync.Once
	file_proto_timeline_timeline_proto_rawDescData = file_proto_timeline_timeline_proto_rawDesc
)

func file_proto_timeline_timeline_proto_rawDescGZIP() []byte {
	file_proto_timeline_timeline_proto_rawDescOnce.Do(func() {
		file_proto_timeline_timeline_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_timeline_timeline_proto_rawDescData)
	})
	return file_proto_timeline_timeline_proto_rawDescData
}

var file_proto_timeline_timeline_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_timeline_timeline_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_timeline_timeline_proto_goTypes = []interface{}{
	(Event)(0),               // 0: Event
	(*TimelineRequest)(nil),  // 1: TimelineRequest
	(*TimelineResponse)(nil), // 2: TimelineResponse
	(*post.Post)(nil),        // 3: Post
}
var file_proto_timeline_timeline_proto_depIdxs = []int32{
	0, // 0: TimelineResponse.event_type:type_name -> Event
	3, // 1: TimelineResponse.posts:type_name -> Post
	1, // 2: TimelineService.GetPosts:input_type -> TimelineRequest
	2, // 3: TimelineService.GetPosts:output_type -> TimelineResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_timeline_timeline_proto_init() }
func file_proto_timeline_timeline_proto_init() {
	if File_proto_timeline_timeline_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_timeline_timeline_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimelineRequest); i {
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
		file_proto_timeline_timeline_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimelineResponse); i {
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
			RawDescriptor: file_proto_timeline_timeline_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_timeline_timeline_proto_goTypes,
		DependencyIndexes: file_proto_timeline_timeline_proto_depIdxs,
		EnumInfos:         file_proto_timeline_timeline_proto_enumTypes,
		MessageInfos:      file_proto_timeline_timeline_proto_msgTypes,
	}.Build()
	File_proto_timeline_timeline_proto = out.File
	file_proto_timeline_timeline_proto_rawDesc = nil
	file_proto_timeline_timeline_proto_goTypes = nil
	file_proto_timeline_timeline_proto_depIdxs = nil
}
