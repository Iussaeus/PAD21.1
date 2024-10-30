// proto/message.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: proto/message.proto

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

type MessageRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender  string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Topic   string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *MessageRequest) Reset() {
	*x = MessageRequest{}
	mi := &file_proto_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageRequest) ProtoMessage() {}

func (x *MessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageRequest.ProtoReflect.Descriptor instead.
func (*MessageRequest) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{0}
}

func (x *MessageRequest) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *MessageRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *MessageRequest) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type MessageResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Confirmation string `protobuf:"bytes,1,opt,name=confirmation,proto3" json:"confirmation,omitempty"`
}

func (x *MessageResponse) Reset() {
	*x = MessageResponse{}
	mi := &file_proto_message_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageResponse) ProtoMessage() {}

func (x *MessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageResponse.ProtoReflect.Descriptor instead.
func (*MessageResponse) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{1}
}

func (x *MessageResponse) GetConfirmation() string {
	if x != nil {
		return x.Confirmation
	}
	return ""
}

type SubscribeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId string `protobuf:"bytes,1,opt,name=clientId,proto3" json:"clientId,omitempty"`
	Topic    string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
}

func (x *SubscribeRequest) Reset() {
	*x = SubscribeRequest{}
	mi := &file_proto_message_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubscribeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubscribeRequest) ProtoMessage() {}

func (x *SubscribeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubscribeRequest.ProtoReflect.Descriptor instead.
func (*SubscribeRequest) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{2}
}

func (x *SubscribeRequest) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *SubscribeRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sender    string `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	Topic     string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Content   string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	Timestamp int64  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	mi := &file_proto_message_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[3]
	if x != nil {
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
	return file_proto_message_proto_rawDescGZIP(), []int{3}
}

func (x *Message) GetSender() string {
	if x != nil {
		return x.Sender
	}
	return ""
}

func (x *Message) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Message) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Message) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type TopicsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TopicsRequest) Reset() {
	*x = TopicsRequest{}
	mi := &file_proto_message_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TopicsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicsRequest) ProtoMessage() {}

func (x *TopicsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicsRequest.ProtoReflect.Descriptor instead.
func (*TopicsRequest) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{4}
}

type TopicsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topics []string `protobuf:"bytes,1,rep,name=topics,proto3" json:"topics,omitempty"`
}

func (x *TopicsResponse) Reset() {
	*x = TopicsResponse{}
	mi := &file_proto_message_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TopicsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicsResponse) ProtoMessage() {}

func (x *TopicsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicsResponse.ProtoReflect.Descriptor instead.
func (*TopicsResponse) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{5}
}

func (x *TopicsResponse) GetTopics() []string {
	if x != nil {
		return x.Topics
	}
	return nil
}

type TopicRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
}

func (x *TopicRequest) Reset() {
	*x = TopicRequest{}
	mi := &file_proto_message_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TopicRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicRequest) ProtoMessage() {}

func (x *TopicRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicRequest.ProtoReflect.Descriptor instead.
func (*TopicRequest) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{6}
}

func (x *TopicRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

type TopicResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *TopicResponse) Reset() {
	*x = TopicResponse{}
	mi := &file_proto_message_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TopicResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopicResponse) ProtoMessage() {}

func (x *TopicResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_message_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopicResponse.ProtoReflect.Descriptor instead.
func (*TopicResponse) Descriptor() ([]byte, []int) {
	return file_proto_message_proto_rawDescGZIP(), []int{7}
}

func (x *TopicResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_proto_message_proto protoreflect.FileDescriptor

var file_proto_message_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x22, 0x58, 0x0a,
	0x0e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x18, 0x0a,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x35, 0x0a, 0x0f, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x44,
	0x0a, 0x10, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x70, 0x69, 0x63, 0x22, 0x6f, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x0f, 0x0a, 0x0d, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x28, 0x0a, 0x0e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x6f, 0x70, 0x69,
	0x63, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x73,
	0x22, 0x24, 0x0a, 0x0c, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x22, 0x27, 0x0a, 0x0d, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32,
	0xbd, 0x02, 0x0a, 0x0d, 0x42, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x3e, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x16, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65,
	0x72, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x38, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62, 0x65, 0x12, 0x18,
	0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x62,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65,
	0x72, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x30, 0x01, 0x12, 0x3a, 0x0a, 0x09, 0x47,
	0x65, 0x74, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x12, 0x15, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65,
	0x72, 0x2e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x14, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e,
	0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x62,
	0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x6f, 0x70,
	0x69, 0x63, 0x12, 0x14, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x2e, 0x54, 0x6f, 0x70, 0x69,
	0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x62, 0x72, 0x6f, 0x6b, 0x65,
	0x72, 0x2e, 0x54, 0x6f, 0x70, 0x69, 0x63, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x09, 0x5a, 0x07, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_proto_message_proto_rawDescOnce sync.Once
	file_proto_message_proto_rawDescData = file_proto_message_proto_rawDesc
)

func file_proto_message_proto_rawDescGZIP() []byte {
	file_proto_message_proto_rawDescOnce.Do(func() {
		file_proto_message_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_message_proto_rawDescData)
	})
	return file_proto_message_proto_rawDescData
}

var file_proto_message_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_proto_message_proto_goTypes = []any{
	(*MessageRequest)(nil),   // 0: broker.MessageRequest
	(*MessageResponse)(nil),  // 1: broker.MessageResponse
	(*SubscribeRequest)(nil), // 2: broker.SubscribeRequest
	(*Message)(nil),          // 3: broker.Message
	(*TopicsRequest)(nil),    // 4: broker.TopicsRequest
	(*TopicsResponse)(nil),   // 5: broker.TopicsResponse
	(*TopicRequest)(nil),     // 6: broker.TopicRequest
	(*TopicResponse)(nil),    // 7: broker.TopicResponse
}
var file_proto_message_proto_depIdxs = []int32{
	0, // 0: broker.BrokerService.SendMessage:input_type -> broker.MessageRequest
	2, // 1: broker.BrokerService.Subscribe:input_type -> broker.SubscribeRequest
	4, // 2: broker.BrokerService.GetTopics:input_type -> broker.TopicsRequest
	6, // 3: broker.BrokerService.CreateTopic:input_type -> broker.TopicRequest
	6, // 4: broker.BrokerService.DeleteTopic:input_type -> broker.TopicRequest
	1, // 5: broker.BrokerService.SendMessage:output_type -> broker.MessageResponse
	3, // 6: broker.BrokerService.Subscribe:output_type -> broker.Message
	5, // 7: broker.BrokerService.GetTopics:output_type -> broker.TopicsResponse
	7, // 8: broker.BrokerService.CreateTopic:output_type -> broker.TopicResponse
	7, // 9: broker.BrokerService.DeleteTopic:output_type -> broker.TopicResponse
	5, // [5:10] is the sub-list for method output_type
	0, // [0:5] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_message_proto_init() }
func file_proto_message_proto_init() {
	if File_proto_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_message_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_message_proto_goTypes,
		DependencyIndexes: file_proto_message_proto_depIdxs,
		MessageInfos:      file_proto_message_proto_msgTypes,
	}.Build()
	File_proto_message_proto = out.File
	file_proto_message_proto_rawDesc = nil
	file_proto_message_proto_goTypes = nil
	file_proto_message_proto_depIdxs = nil
}