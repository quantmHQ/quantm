// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: hooks/github/v1/github.proto

package githubv1

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

// Represents a provider's repo within the control plane.
type GithubRepo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CreatedAt      *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt      *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	RepoId         string                 `protobuf:"bytes,4,opt,name=repo_id,json=repoId,proto3" json:"repo_id,omitempty"` // Refers to the core repo
	InstallationId string                 `protobuf:"bytes,5,opt,name=installation_id,json=installationId,proto3" json:"installation_id,omitempty"`
	GithubId       int64                  `protobuf:"varint,6,opt,name=github_id,json=githubId,proto3" json:"github_id,omitempty"`
	Name           string                 `protobuf:"bytes,7,opt,name=name,proto3" json:"name,omitempty"`
	FullName       string                 `protobuf:"bytes,8,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Url            string                 `protobuf:"bytes,9,opt,name=url,proto3" json:"url,omitempty"`
	IsActive       bool                   `protobuf:"varint,10,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
}

func (x *GithubRepo) Reset() {
	*x = GithubRepo{}
	mi := &file_hooks_github_v1_github_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GithubRepo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GithubRepo) ProtoMessage() {}

func (x *GithubRepo) ProtoReflect() protoreflect.Message {
	mi := &file_hooks_github_v1_github_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GithubRepo.ProtoReflect.Descriptor instead.
func (*GithubRepo) Descriptor() ([]byte, []int) {
	return file_hooks_github_v1_github_proto_rawDescGZIP(), []int{0}
}

func (x *GithubRepo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GithubRepo) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *GithubRepo) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *GithubRepo) GetRepoId() string {
	if x != nil {
		return x.RepoId
	}
	return ""
}

func (x *GithubRepo) GetInstallationId() string {
	if x != nil {
		return x.InstallationId
	}
	return ""
}

func (x *GithubRepo) GetGithubId() int64 {
	if x != nil {
		return x.GithubId
	}
	return 0
}

func (x *GithubRepo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GithubRepo) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *GithubRepo) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *GithubRepo) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

// Request to create a org's provider's repo.
type CreateGithubRepoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name           string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	FullName       string `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Url            string `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	IsActive       bool   `protobuf:"varint,4,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
	RepoId         string `protobuf:"bytes,5,opt,name=repo_id,json=repoId,proto3" json:"repo_id,omitempty"`
	InstallationId string `protobuf:"bytes,6,opt,name=installation_id,json=installationId,proto3" json:"installation_id,omitempty"`
}

func (x *CreateGithubRepoRequest) Reset() {
	*x = CreateGithubRepoRequest{}
	mi := &file_hooks_github_v1_github_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateGithubRepoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateGithubRepoRequest) ProtoMessage() {}

func (x *CreateGithubRepoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_hooks_github_v1_github_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateGithubRepoRequest.ProtoReflect.Descriptor instead.
func (*CreateGithubRepoRequest) Descriptor() ([]byte, []int) {
	return file_hooks_github_v1_github_proto_rawDescGZIP(), []int{1}
}

func (x *CreateGithubRepoRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateGithubRepoRequest) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *CreateGithubRepoRequest) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *CreateGithubRepoRequest) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

func (x *CreateGithubRepoRequest) GetRepoId() string {
	if x != nil {
		return x.RepoId
	}
	return ""
}

func (x *CreateGithubRepoRequest) GetInstallationId() string {
	if x != nil {
		return x.InstallationId
	}
	return ""
}

// Response to get org's provider's repo.
type CreateGithubRepoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GithubRepo *GithubRepo `protobuf:"bytes,1,opt,name=github_repo,json=githubRepo,proto3" json:"github_repo,omitempty"`
}

func (x *CreateGithubRepoResponse) Reset() {
	*x = CreateGithubRepoResponse{}
	mi := &file_hooks_github_v1_github_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateGithubRepoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateGithubRepoResponse) ProtoMessage() {}

func (x *CreateGithubRepoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_hooks_github_v1_github_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateGithubRepoResponse.ProtoReflect.Descriptor instead.
func (*CreateGithubRepoResponse) Descriptor() ([]byte, []int) {
	return file_hooks_github_v1_github_proto_rawDescGZIP(), []int{2}
}

func (x *CreateGithubRepoResponse) GetGithubRepo() *GithubRepo {
	if x != nil {
		return x.GithubRepo
	}
	return nil
}

// Request to get org's provider's repo by id.
type GetGithubRepoByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetGithubRepoByIDRequest) Reset() {
	*x = GetGithubRepoByIDRequest{}
	mi := &file_hooks_github_v1_github_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetGithubRepoByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGithubRepoByIDRequest) ProtoMessage() {}

func (x *GetGithubRepoByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_hooks_github_v1_github_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGithubRepoByIDRequest.ProtoReflect.Descriptor instead.
func (*GetGithubRepoByIDRequest) Descriptor() ([]byte, []int) {
	return file_hooks_github_v1_github_proto_rawDescGZIP(), []int{3}
}

func (x *GetGithubRepoByIDRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Response to get org's provider's repo.
type GetGithubRepoByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GithubRepo *GithubRepo `protobuf:"bytes,1,opt,name=github_repo,json=githubRepo,proto3" json:"github_repo,omitempty"`
}

func (x *GetGithubRepoByIDResponse) Reset() {
	*x = GetGithubRepoByIDResponse{}
	mi := &file_hooks_github_v1_github_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetGithubRepoByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGithubRepoByIDResponse) ProtoMessage() {}

func (x *GetGithubRepoByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_hooks_github_v1_github_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGithubRepoByIDResponse.ProtoReflect.Descriptor instead.
func (*GetGithubRepoByIDResponse) Descriptor() ([]byte, []int) {
	return file_hooks_github_v1_github_proto_rawDescGZIP(), []int{4}
}

func (x *GetGithubRepoByIDResponse) GetGithubRepo() *GithubRepo {
	if x != nil {
		return x.GithubRepo
	}
	return nil
}

// Request to get org's provider's repo by name.
type GetGithubRepoByNameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetGithubRepoByNameRequest) Reset() {
	*x = GetGithubRepoByNameRequest{}
	mi := &file_hooks_github_v1_github_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetGithubRepoByNameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGithubRepoByNameRequest) ProtoMessage() {}

func (x *GetGithubRepoByNameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_hooks_github_v1_github_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGithubRepoByNameRequest.ProtoReflect.Descriptor instead.
func (*GetGithubRepoByNameRequest) Descriptor() ([]byte, []int) {
	return file_hooks_github_v1_github_proto_rawDescGZIP(), []int{5}
}

func (x *GetGithubRepoByNameRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

// Response to get org's provider's repo.
type GetGithubRepoByNameResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GithubRepo *GithubRepo `protobuf:"bytes,1,opt,name=github_repo,json=githubRepo,proto3" json:"github_repo,omitempty"`
}

func (x *GetGithubRepoByNameResponse) Reset() {
	*x = GetGithubRepoByNameResponse{}
	mi := &file_hooks_github_v1_github_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetGithubRepoByNameResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGithubRepoByNameResponse) ProtoMessage() {}

func (x *GetGithubRepoByNameResponse) ProtoReflect() protoreflect.Message {
	mi := &file_hooks_github_v1_github_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGithubRepoByNameResponse.ProtoReflect.Descriptor instead.
func (*GetGithubRepoByNameResponse) Descriptor() ([]byte, []int) {
	return file_hooks_github_v1_github_proto_rawDescGZIP(), []int{6}
}

func (x *GetGithubRepoByNameResponse) GetGithubRepo() *GithubRepo {
	if x != nil {
		return x.GithubRepo
	}
	return nil
}

var File_hooks_github_v1_github_proto protoreflect.FileDescriptor

var file_hooks_github_v1_github_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2f, 0x76,
	0x31, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f,
	0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xd1, 0x02, 0x0a, 0x0a, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x72, 0x65, 0x70, 0x6f, 0x5f, 0x69, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x49, 0x64, 0x12, 0x27,
	0x0a, 0x0f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x75, 0x6c, 0x6c,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c,
	0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x61, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x41, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x22, 0xbb, 0x01, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x75, 0x6c, 0x6c, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x75, 0x6c, 0x6c, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x75, 0x72, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x12, 0x17, 0x0a, 0x07, 0x72, 0x65, 0x70, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x49, 0x64, 0x12, 0x27, 0x0a, 0x0f, 0x69, 0x6e, 0x73,
	0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x22, 0x58, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c,
	0x0a, 0x0b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f,
	0x52, 0x0a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x22, 0x2a, 0x0a, 0x18,
	0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x49,
	0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x59, 0x0a, 0x19, 0x47, 0x65, 0x74, 0x47,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x0b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x5f,
	0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x68, 0x6f, 0x6f,
	0x6b, 0x73, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x0a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52,
	0x65, 0x70, 0x6f, 0x22, 0x30, 0x0a, 0x1a, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x5b, 0x0a, 0x1b, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x0b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x5f, 0x72,
	0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x68, 0x6f, 0x6f, 0x6b,
	0x73, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x0a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65,
	0x70, 0x6f, 0x32, 0xce, 0x02, 0x0a, 0x0b, 0x52, 0x65, 0x70, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x61, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x12, 0x28, 0x2e, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52,
	0x65, 0x70, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x68, 0x6f, 0x6f,
	0x6b, 0x73, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6a, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x49, 0x44, 0x12, 0x29, 0x2e, 0x68, 0x6f, 0x6f,
	0x6b, 0x73, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x70, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x52, 0x65,
	0x70, 0x6f, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x2b, 0x2e, 0x68, 0x6f, 0x6f, 0x6b, 0x73,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x47, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x68, 0x6f, 0x6f, 0x6b, 0x73, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x47, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x52, 0x65, 0x70, 0x6f, 0x42, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0xbb, 0x01, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x6f, 0x6f, 0x6b,
	0x73, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x76, 0x31, 0x42, 0x0b, 0x47, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x6f, 0x2e, 0x62,
	0x72, 0x65, 0x75, 0x2e, 0x69, 0x6f, 0x2f, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x6d, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x68, 0x6f, 0x6f,
	0x6b, 0x73, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2f, 0x76, 0x31, 0x3b, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x48, 0x47, 0x58, 0xaa, 0x02, 0x0f, 0x48, 0x6f,
	0x6f, 0x6b, 0x73, 0x2e, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0f,
	0x48, 0x6f, 0x6f, 0x6b, 0x73, 0x5c, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x5c, 0x56, 0x31, 0xe2,
	0x02, 0x1b, 0x48, 0x6f, 0x6f, 0x6b, 0x73, 0x5c, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x5c, 0x56,
	0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x11,
	0x48, 0x6f, 0x6f, 0x6b, 0x73, 0x3a, 0x3a, 0x47, 0x69, 0x74, 0x68, 0x75, 0x62, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_hooks_github_v1_github_proto_rawDescOnce sync.Once
	file_hooks_github_v1_github_proto_rawDescData = file_hooks_github_v1_github_proto_rawDesc
)

func file_hooks_github_v1_github_proto_rawDescGZIP() []byte {
	file_hooks_github_v1_github_proto_rawDescOnce.Do(func() {
		file_hooks_github_v1_github_proto_rawDescData = protoimpl.X.CompressGZIP(file_hooks_github_v1_github_proto_rawDescData)
	})
	return file_hooks_github_v1_github_proto_rawDescData
}

var file_hooks_github_v1_github_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_hooks_github_v1_github_proto_goTypes = []any{
	(*GithubRepo)(nil),                  // 0: hooks.github.v1.GithubRepo
	(*CreateGithubRepoRequest)(nil),     // 1: hooks.github.v1.CreateGithubRepoRequest
	(*CreateGithubRepoResponse)(nil),    // 2: hooks.github.v1.CreateGithubRepoResponse
	(*GetGithubRepoByIDRequest)(nil),    // 3: hooks.github.v1.GetGithubRepoByIDRequest
	(*GetGithubRepoByIDResponse)(nil),   // 4: hooks.github.v1.GetGithubRepoByIDResponse
	(*GetGithubRepoByNameRequest)(nil),  // 5: hooks.github.v1.GetGithubRepoByNameRequest
	(*GetGithubRepoByNameResponse)(nil), // 6: hooks.github.v1.GetGithubRepoByNameResponse
	(*timestamppb.Timestamp)(nil),       // 7: google.protobuf.Timestamp
}
var file_hooks_github_v1_github_proto_depIdxs = []int32{
	7, // 0: hooks.github.v1.GithubRepo.created_at:type_name -> google.protobuf.Timestamp
	7, // 1: hooks.github.v1.GithubRepo.updated_at:type_name -> google.protobuf.Timestamp
	0, // 2: hooks.github.v1.CreateGithubRepoResponse.github_repo:type_name -> hooks.github.v1.GithubRepo
	0, // 3: hooks.github.v1.GetGithubRepoByIDResponse.github_repo:type_name -> hooks.github.v1.GithubRepo
	0, // 4: hooks.github.v1.GetGithubRepoByNameResponse.github_repo:type_name -> hooks.github.v1.GithubRepo
	1, // 5: hooks.github.v1.RepoService.CreateRepo:input_type -> hooks.github.v1.CreateGithubRepoRequest
	3, // 6: hooks.github.v1.RepoService.GetGithubRepoByID:input_type -> hooks.github.v1.GetGithubRepoByIDRequest
	5, // 7: hooks.github.v1.RepoService.GetGithubRepoByName:input_type -> hooks.github.v1.GetGithubRepoByNameRequest
	2, // 8: hooks.github.v1.RepoService.CreateRepo:output_type -> hooks.github.v1.CreateGithubRepoResponse
	4, // 9: hooks.github.v1.RepoService.GetGithubRepoByID:output_type -> hooks.github.v1.GetGithubRepoByIDResponse
	6, // 10: hooks.github.v1.RepoService.GetGithubRepoByName:output_type -> hooks.github.v1.GetGithubRepoByNameResponse
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_hooks_github_v1_github_proto_init() }
func file_hooks_github_v1_github_proto_init() {
	if File_hooks_github_v1_github_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_hooks_github_v1_github_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_hooks_github_v1_github_proto_goTypes,
		DependencyIndexes: file_hooks_github_v1_github_proto_depIdxs,
		MessageInfos:      file_hooks_github_v1_github_proto_msgTypes,
	}.Build()
	File_hooks_github_v1_github_proto = out.File
	file_hooks_github_v1_github_proto_rawDesc = nil
	file_hooks_github_v1_github_proto_goTypes = nil
	file_hooks_github_v1_github_proto_depIdxs = nil
}
