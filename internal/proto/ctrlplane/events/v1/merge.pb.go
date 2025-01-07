// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: ctrlplane/events/v1/merge.proto

package eventsv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Merge struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	HeadBranch    string                 `protobuf:"bytes,1,opt,name=head_branch,json=headBranch,proto3" json:"head_branch,omitempty"`
	HeadCommit    *Commit                `protobuf:"bytes,2,opt,name=head_commit,json=headCommit,proto3" json:"head_commit,omitempty"`
	BaseBranch    string                 `protobuf:"bytes,3,opt,name=base_branch,json=baseBranch,proto3" json:"base_branch,omitempty"`
	BaseCommit    *Commit                `protobuf:"bytes,4,opt,name=base_commit,json=baseCommit,proto3" json:"base_commit,omitempty"`
	Files         []string               `protobuf:"bytes,5,rep,name=files,proto3" json:"files,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Merge) Reset() {
	*x = Merge{}
	mi := &file_ctrlplane_events_v1_merge_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Merge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Merge) ProtoMessage() {}

func (x *Merge) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_events_v1_merge_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Merge.ProtoReflect.Descriptor instead.
func (*Merge) Descriptor() ([]byte, []int) {
	return file_ctrlplane_events_v1_merge_proto_rawDescGZIP(), []int{0}
}

func (x *Merge) GetHeadBranch() string {
	if x != nil {
		return x.HeadBranch
	}
	return ""
}

func (x *Merge) GetHeadCommit() *Commit {
	if x != nil {
		return x.HeadCommit
	}
	return nil
}

func (x *Merge) GetBaseBranch() string {
	if x != nil {
		return x.BaseBranch
	}
	return ""
}

func (x *Merge) GetBaseCommit() *Commit {
	if x != nil {
		return x.BaseCommit
	}
	return nil
}

func (x *Merge) GetFiles() []string {
	if x != nil {
		return x.Files
	}
	return nil
}

type MergeQueue struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Number        int64                  `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
	Branch        string                 `protobuf:"bytes,2,opt,name=branch,proto3" json:"branch,omitempty"`
	IsPriority    bool                   `protobuf:"varint,3,opt,name=is_priority,json=isPriority,proto3" json:"is_priority,omitempty"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MergeQueue) Reset() {
	*x = MergeQueue{}
	mi := &file_ctrlplane_events_v1_merge_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MergeQueue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MergeQueue) ProtoMessage() {}

func (x *MergeQueue) ProtoReflect() protoreflect.Message {
	mi := &file_ctrlplane_events_v1_merge_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MergeQueue.ProtoReflect.Descriptor instead.
func (*MergeQueue) Descriptor() ([]byte, []int) {
	return file_ctrlplane_events_v1_merge_proto_rawDescGZIP(), []int{1}
}

func (x *MergeQueue) GetNumber() int64 {
	if x != nil {
		return x.Number
	}
	return 0
}

func (x *MergeQueue) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *MergeQueue) GetIsPriority() bool {
	if x != nil {
		return x.IsPriority
	}
	return false
}

func (x *MergeQueue) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

var File_ctrlplane_events_v1_merge_proto protoreflect.FileDescriptor

var file_ctrlplane_events_v1_merge_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x13, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x20, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e,
	0x65, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xdb, 0x01, 0x0a, 0x05, 0x4d, 0x65,
	0x72, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x65, 0x61, 0x64, 0x5f, 0x62, 0x72, 0x61, 0x6e,
	0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x68, 0x65, 0x61, 0x64, 0x42, 0x72,
	0x61, 0x6e, 0x63, 0x68, 0x12, 0x3c, 0x0a, 0x0b, 0x68, 0x65, 0x61, 0x64, 0x5f, 0x63, 0x6f, 0x6d,
	0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x74, 0x72, 0x6c,
	0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x0a, 0x68, 0x65, 0x61, 0x64, 0x43, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x62, 0x72, 0x61, 0x6e, 0x63,
	0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x42, 0x72, 0x61,
	0x6e, 0x63, 0x68, 0x12, 0x3c, 0x0a, 0x0b, 0x62, 0x61, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d,
	0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70,
	0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x52, 0x0a, 0x62, 0x61, 0x73, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x69,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x97, 0x01, 0x0a, 0x0a, 0x4d, 0x65, 0x72, 0x67,
	0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x62, 0x72, 0x61, 0x6e, 0x63, 0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x69, 0x73, 0x5f, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x69, 0x73, 0x50,
	0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x42, 0xd2, 0x01, 0x0a, 0x17, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x74, 0x72, 0x6c, 0x70, 0x6c,
	0x61, 0x6e, 0x65, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x4d,
	0x65, 0x72, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3d, 0x67, 0x6f, 0x2e,
	0x62, 0x72, 0x65, 0x75, 0x2e, 0x69, 0x6f, 0x2f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x6d, 0x2f, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x74,
	0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x76,
	0x31, 0x3b, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x45, 0x58,
	0xaa, 0x02, 0x13, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x2e, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x73, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x13, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61,
	0x6e, 0x65, 0x5c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1f, 0x43,
	0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x5c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x5c,
	0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x15, 0x43, 0x74, 0x72, 0x6c, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x3a, 0x3a, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ctrlplane_events_v1_merge_proto_rawDescOnce sync.Once
	file_ctrlplane_events_v1_merge_proto_rawDescData = file_ctrlplane_events_v1_merge_proto_rawDesc
)

func file_ctrlplane_events_v1_merge_proto_rawDescGZIP() []byte {
	file_ctrlplane_events_v1_merge_proto_rawDescOnce.Do(func() {
		file_ctrlplane_events_v1_merge_proto_rawDescData = protoimpl.X.CompressGZIP(file_ctrlplane_events_v1_merge_proto_rawDescData)
	})
	return file_ctrlplane_events_v1_merge_proto_rawDescData
}

var file_ctrlplane_events_v1_merge_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_ctrlplane_events_v1_merge_proto_goTypes = []any{
	(*Merge)(nil),                 // 0: ctrlplane.events.v1.Merge
	(*MergeQueue)(nil),            // 1: ctrlplane.events.v1.MergeQueue
	(*Commit)(nil),                // 2: ctrlplane.events.v1.Commit
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_ctrlplane_events_v1_merge_proto_depIdxs = []int32{
	2, // 0: ctrlplane.events.v1.Merge.head_commit:type_name -> ctrlplane.events.v1.Commit
	2, // 1: ctrlplane.events.v1.Merge.base_commit:type_name -> ctrlplane.events.v1.Commit
	3, // 2: ctrlplane.events.v1.MergeQueue.timestamp:type_name -> google.protobuf.Timestamp
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_ctrlplane_events_v1_merge_proto_init() }
func file_ctrlplane_events_v1_merge_proto_init() {
	if File_ctrlplane_events_v1_merge_proto != nil {
		return
	}
	file_ctrlplane_events_v1_commit_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ctrlplane_events_v1_merge_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ctrlplane_events_v1_merge_proto_goTypes,
		DependencyIndexes: file_ctrlplane_events_v1_merge_proto_depIdxs,
		MessageInfos:      file_ctrlplane_events_v1_merge_proto_msgTypes,
	}.Build()
	File_ctrlplane_events_v1_merge_proto = out.File
	file_ctrlplane_events_v1_merge_proto_rawDesc = nil
	file_ctrlplane_events_v1_merge_proto_goTypes = nil
	file_ctrlplane_events_v1_merge_proto_depIdxs = nil
}
