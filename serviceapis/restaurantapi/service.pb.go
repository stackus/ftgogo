// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.6
// source: service.proto

package restaurantapi

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

type CreateRestaurantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string                          `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Address *CreateRestaurantRequestAddress `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
	Menu    *CreateRestaurantRequestMenu    `protobuf:"bytes,3,opt,name=Menu,proto3" json:"Menu,omitempty"`
}

func (x *CreateRestaurantRequest) Reset() {
	*x = CreateRestaurantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRestaurantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRestaurantRequest) ProtoMessage() {}

func (x *CreateRestaurantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRestaurantRequest.ProtoReflect.Descriptor instead.
func (*CreateRestaurantRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateRestaurantRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateRestaurantRequest) GetAddress() *CreateRestaurantRequestAddress {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *CreateRestaurantRequest) GetMenu() *CreateRestaurantRequestMenu {
	if x != nil {
		return x.Menu
	}
	return nil
}

type CreateRestaurantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RestaurantID string `protobuf:"bytes,1,opt,name=RestaurantID,proto3" json:"RestaurantID,omitempty"`
}

func (x *CreateRestaurantResponse) Reset() {
	*x = CreateRestaurantResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRestaurantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRestaurantResponse) ProtoMessage() {}

func (x *CreateRestaurantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRestaurantResponse.ProtoReflect.Descriptor instead.
func (*CreateRestaurantResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateRestaurantResponse) GetRestaurantID() string {
	if x != nil {
		return x.RestaurantID
	}
	return ""
}

type GetRestaurantRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RestaurantID string `protobuf:"bytes,1,opt,name=RestaurantID,proto3" json:"RestaurantID,omitempty"`
}

func (x *GetRestaurantRequest) Reset() {
	*x = GetRestaurantRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRestaurantRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRestaurantRequest) ProtoMessage() {}

func (x *GetRestaurantRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRestaurantRequest.ProtoReflect.Descriptor instead.
func (*GetRestaurantRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetRestaurantRequest) GetRestaurantID() string {
	if x != nil {
		return x.RestaurantID
	}
	return ""
}

type GetRestaurantResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RestaurantID string                        `protobuf:"bytes,1,opt,name=RestaurantID,proto3" json:"RestaurantID,omitempty"`
	Name         string                        `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Address      *GetRestaurantResponseAddress `protobuf:"bytes,3,opt,name=Address,proto3" json:"Address,omitempty"`
	Menu         *GetRestaurantResponseMenu    `protobuf:"bytes,4,opt,name=Menu,proto3" json:"Menu,omitempty"`
}

func (x *GetRestaurantResponse) Reset() {
	*x = GetRestaurantResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRestaurantResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRestaurantResponse) ProtoMessage() {}

func (x *GetRestaurantResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRestaurantResponse.ProtoReflect.Descriptor instead.
func (*GetRestaurantResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetRestaurantResponse) GetRestaurantID() string {
	if x != nil {
		return x.RestaurantID
	}
	return ""
}

func (x *GetRestaurantResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetRestaurantResponse) GetAddress() *GetRestaurantResponseAddress {
	if x != nil {
		return x.Address
	}
	return nil
}

func (x *GetRestaurantResponse) GetMenu() *GetRestaurantResponseMenu {
	if x != nil {
		return x.Menu
	}
	return nil
}

type CreateRestaurantRequestAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Street1 string `protobuf:"bytes,1,opt,name=Street1,proto3" json:"Street1,omitempty"`
	Street2 string `protobuf:"bytes,2,opt,name=Street2,proto3" json:"Street2,omitempty"`
	City    string `protobuf:"bytes,3,opt,name=City,proto3" json:"City,omitempty"`
	State   string `protobuf:"bytes,4,opt,name=State,proto3" json:"State,omitempty"`
	Zip     string `protobuf:"bytes,5,opt,name=Zip,proto3" json:"Zip,omitempty"`
}

func (x *CreateRestaurantRequestAddress) Reset() {
	*x = CreateRestaurantRequestAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRestaurantRequestAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRestaurantRequestAddress) ProtoMessage() {}

func (x *CreateRestaurantRequestAddress) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRestaurantRequestAddress.ProtoReflect.Descriptor instead.
func (*CreateRestaurantRequestAddress) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0, 0}
}

func (x *CreateRestaurantRequestAddress) GetStreet1() string {
	if x != nil {
		return x.Street1
	}
	return ""
}

func (x *CreateRestaurantRequestAddress) GetStreet2() string {
	if x != nil {
		return x.Street2
	}
	return ""
}

func (x *CreateRestaurantRequestAddress) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *CreateRestaurantRequestAddress) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *CreateRestaurantRequestAddress) GetZip() string {
	if x != nil {
		return x.Zip
	}
	return ""
}

type CreateRestaurantRequestMenuItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID    string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Price int64  `protobuf:"varint,3,opt,name=Price,proto3" json:"Price,omitempty"`
}

func (x *CreateRestaurantRequestMenuItem) Reset() {
	*x = CreateRestaurantRequestMenuItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRestaurantRequestMenuItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRestaurantRequestMenuItem) ProtoMessage() {}

func (x *CreateRestaurantRequestMenuItem) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRestaurantRequestMenuItem.ProtoReflect.Descriptor instead.
func (*CreateRestaurantRequestMenuItem) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0, 1}
}

func (x *CreateRestaurantRequestMenuItem) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *CreateRestaurantRequestMenuItem) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateRestaurantRequestMenuItem) GetPrice() int64 {
	if x != nil {
		return x.Price
	}
	return 0
}

type CreateRestaurantRequestMenu struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MenuItems []*CreateRestaurantRequestMenuItem `protobuf:"bytes,1,rep,name=MenuItems,proto3" json:"MenuItems,omitempty"`
}

func (x *CreateRestaurantRequestMenu) Reset() {
	*x = CreateRestaurantRequestMenu{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRestaurantRequestMenu) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRestaurantRequestMenu) ProtoMessage() {}

func (x *CreateRestaurantRequestMenu) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRestaurantRequestMenu.ProtoReflect.Descriptor instead.
func (*CreateRestaurantRequestMenu) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0, 2}
}

func (x *CreateRestaurantRequestMenu) GetMenuItems() []*CreateRestaurantRequestMenuItem {
	if x != nil {
		return x.MenuItems
	}
	return nil
}

type GetRestaurantResponseAddress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Street1 string `protobuf:"bytes,1,opt,name=Street1,proto3" json:"Street1,omitempty"`
	Street2 string `protobuf:"bytes,2,opt,name=Street2,proto3" json:"Street2,omitempty"`
	City    string `protobuf:"bytes,3,opt,name=City,proto3" json:"City,omitempty"`
	State   string `protobuf:"bytes,4,opt,name=State,proto3" json:"State,omitempty"`
	Zip     string `protobuf:"bytes,5,opt,name=Zip,proto3" json:"Zip,omitempty"`
}

func (x *GetRestaurantResponseAddress) Reset() {
	*x = GetRestaurantResponseAddress{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRestaurantResponseAddress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRestaurantResponseAddress) ProtoMessage() {}

func (x *GetRestaurantResponseAddress) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRestaurantResponseAddress.ProtoReflect.Descriptor instead.
func (*GetRestaurantResponseAddress) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3, 0}
}

func (x *GetRestaurantResponseAddress) GetStreet1() string {
	if x != nil {
		return x.Street1
	}
	return ""
}

func (x *GetRestaurantResponseAddress) GetStreet2() string {
	if x != nil {
		return x.Street2
	}
	return ""
}

func (x *GetRestaurantResponseAddress) GetCity() string {
	if x != nil {
		return x.City
	}
	return ""
}

func (x *GetRestaurantResponseAddress) GetState() string {
	if x != nil {
		return x.State
	}
	return ""
}

func (x *GetRestaurantResponseAddress) GetZip() string {
	if x != nil {
		return x.Zip
	}
	return ""
}

type GetRestaurantResponseMenuItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID    string `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Name  string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	Price int64  `protobuf:"varint,3,opt,name=Price,proto3" json:"Price,omitempty"`
}

func (x *GetRestaurantResponseMenuItem) Reset() {
	*x = GetRestaurantResponseMenuItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRestaurantResponseMenuItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRestaurantResponseMenuItem) ProtoMessage() {}

func (x *GetRestaurantResponseMenuItem) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRestaurantResponseMenuItem.ProtoReflect.Descriptor instead.
func (*GetRestaurantResponseMenuItem) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3, 1}
}

func (x *GetRestaurantResponseMenuItem) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *GetRestaurantResponseMenuItem) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetRestaurantResponseMenuItem) GetPrice() int64 {
	if x != nil {
		return x.Price
	}
	return 0
}

type GetRestaurantResponseMenu struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MenuItems []*GetRestaurantResponseMenuItem `protobuf:"bytes,1,rep,name=MenuItems,proto3" json:"MenuItems,omitempty"`
}

func (x *GetRestaurantResponseMenu) Reset() {
	*x = GetRestaurantResponseMenu{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRestaurantResponseMenu) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRestaurantResponseMenu) ProtoMessage() {}

func (x *GetRestaurantResponseMenu) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRestaurantResponseMenu.ProtoReflect.Descriptor instead.
func (*GetRestaurantResponseMenu) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3, 2}
}

func (x *GetRestaurantResponseMenu) GetMenuItems() []*GetRestaurantResponseMenuItem {
	if x != nil {
		return x.MenuItems
	}
	return nil
}

var File_service_proto protoreflect.FileDescriptor

var file_service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0d, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x61, 0x70, 0x69, 0x22, 0xd0,
	0x03, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72,
	0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x48,
	0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2e, 0x2e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x61, 0x70, 0x69, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52,
	0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3f, 0x0a, 0x04, 0x4d, 0x65, 0x6e, 0x75,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b, 0x2e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72,
	0x61, 0x6e, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x6d,
	0x65, 0x6e, 0x75, 0x52, 0x04, 0x4d, 0x65, 0x6e, 0x75, 0x1a, 0x79, 0x0a, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x31, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x31, 0x12, 0x18,
	0x0a, 0x07, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x32, 0x12, 0x12, 0x0a, 0x04, 0x43, 0x69, 0x74, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x43, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x5a, 0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x5a, 0x69, 0x70, 0x1a, 0x44, 0x0a, 0x08, 0x6d, 0x65, 0x6e, 0x75, 0x49, 0x74, 0x65, 0x6d,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x12, 0x0a, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x1a, 0x55, 0x0a, 0x04, 0x6d, 0x65,
	0x6e, 0x75, 0x12, 0x4d, 0x0a, 0x09, 0x4d, 0x65, 0x6e, 0x75, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61,
	0x6e, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x74,
	0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x6d, 0x65,
	0x6e, 0x75, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x09, 0x4d, 0x65, 0x6e, 0x75, 0x49, 0x74, 0x65, 0x6d,
	0x73, 0x22, 0x3e, 0x0a, 0x18, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x74, 0x61,
	0x75, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a,
	0x0c, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0c, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x49,
	0x44, 0x22, 0x3a, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x52, 0x65, 0x73,
	0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x49, 0x44, 0x22, 0xec, 0x03,
	0x0a, 0x15, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x74, 0x61,
	0x75, 0x72, 0x61, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x52,
	0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x46, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2c, 0x2e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x61, 0x70, 0x69,
	0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x07,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3d, 0x0a, 0x04, 0x4d, 0x65, 0x6e, 0x75, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61,
	0x6e, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72,
	0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x6d, 0x65, 0x6e, 0x75,
	0x52, 0x04, 0x4d, 0x65, 0x6e, 0x75, 0x1a, 0x79, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x31, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x53, 0x74, 0x72, 0x65, 0x65, 0x74, 0x31, 0x12, 0x18, 0x0a, 0x07, 0x53,
	0x74, 0x72, 0x65, 0x65, 0x74, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x74,
	0x72, 0x65, 0x65, 0x74, 0x32, 0x12, 0x12, 0x0a, 0x04, 0x43, 0x69, 0x74, 0x79, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x43, 0x69, 0x74, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x10, 0x0a, 0x03, 0x5a, 0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x5a, 0x69,
	0x70, 0x1a, 0x44, 0x0a, 0x08, 0x6d, 0x65, 0x6e, 0x75, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a,
	0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x12, 0x0a,
	0x04, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x1a, 0x53, 0x0a, 0x04, 0x6d, 0x65, 0x6e, 0x75, 0x12,
	0x4b, 0x0a, 0x09, 0x4d, 0x65, 0x6e, 0x75, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2d, 0x2e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x61,
	0x70, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x6d, 0x65, 0x6e, 0x75, 0x49, 0x74, 0x65,
	0x6d, 0x52, 0x09, 0x4d, 0x65, 0x6e, 0x75, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x32, 0xd4, 0x01, 0x0a,
	0x11, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x63, 0x0a, 0x10, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x74,
	0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x12, 0x26, 0x2e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72,
	0x61, 0x6e, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27,
	0x2e, 0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x5a, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x12, 0x23, 0x2e, 0x72, 0x65, 0x73, 0x74, 0x61,
	0x75, 0x72, 0x61, 0x6e, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x73, 0x74,
	0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e,
	0x72, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x61, 0x70, 0x69, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x75, 0x73, 0x2f, 0x66, 0x74, 0x67, 0x6f, 0x67, 0x6f,
	0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x72, 0x65, 0x73,
	0x74, 0x61, 0x75, 0x72, 0x61, 0x6e, 0x74, 0x61, 0x70, 0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_service_proto_rawDescOnce sync.Once
	file_service_proto_rawDescData = file_service_proto_rawDesc
)

func file_service_proto_rawDescGZIP() []byte {
	file_service_proto_rawDescOnce.Do(func() {
		file_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_proto_rawDescData)
	})
	return file_service_proto_rawDescData
}

var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_service_proto_goTypes = []interface{}{
	(*CreateRestaurantRequest)(nil),         // 0: restaurantapi.CreateRestaurantRequest
	(*CreateRestaurantResponse)(nil),        // 1: restaurantapi.CreateRestaurantResponse
	(*GetRestaurantRequest)(nil),            // 2: restaurantapi.GetRestaurantRequest
	(*GetRestaurantResponse)(nil),           // 3: restaurantapi.GetRestaurantResponse
	(*CreateRestaurantRequestAddress)(nil),  // 4: restaurantapi.CreateRestaurantRequest.address
	(*CreateRestaurantRequestMenuItem)(nil), // 5: restaurantapi.CreateRestaurantRequest.menuItem
	(*CreateRestaurantRequestMenu)(nil),     // 6: restaurantapi.CreateRestaurantRequest.menu
	(*GetRestaurantResponseAddress)(nil),    // 7: restaurantapi.GetRestaurantResponse.address
	(*GetRestaurantResponseMenuItem)(nil),   // 8: restaurantapi.GetRestaurantResponse.menuItem
	(*GetRestaurantResponseMenu)(nil),       // 9: restaurantapi.GetRestaurantResponse.menu
}
var file_service_proto_depIdxs = []int32{
	4, // 0: restaurantapi.CreateRestaurantRequest.Address:type_name -> restaurantapi.CreateRestaurantRequest.address
	6, // 1: restaurantapi.CreateRestaurantRequest.Menu:type_name -> restaurantapi.CreateRestaurantRequest.menu
	7, // 2: restaurantapi.GetRestaurantResponse.Address:type_name -> restaurantapi.GetRestaurantResponse.address
	9, // 3: restaurantapi.GetRestaurantResponse.Menu:type_name -> restaurantapi.GetRestaurantResponse.menu
	5, // 4: restaurantapi.CreateRestaurantRequest.menu.MenuItems:type_name -> restaurantapi.CreateRestaurantRequest.menuItem
	8, // 5: restaurantapi.GetRestaurantResponse.menu.MenuItems:type_name -> restaurantapi.GetRestaurantResponse.menuItem
	0, // 6: restaurantapi.RestaurantService.CreateRestaurant:input_type -> restaurantapi.CreateRestaurantRequest
	2, // 7: restaurantapi.RestaurantService.GetRestaurant:input_type -> restaurantapi.GetRestaurantRequest
	1, // 8: restaurantapi.RestaurantService.CreateRestaurant:output_type -> restaurantapi.CreateRestaurantResponse
	3, // 9: restaurantapi.RestaurantService.GetRestaurant:output_type -> restaurantapi.GetRestaurantResponse
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRestaurantRequest); i {
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
		file_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRestaurantResponse); i {
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
		file_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRestaurantRequest); i {
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
		file_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRestaurantResponse); i {
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
		file_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRestaurantRequestAddress); i {
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
		file_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRestaurantRequestMenuItem); i {
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
		file_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRestaurantRequestMenu); i {
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
		file_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRestaurantResponseAddress); i {
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
		file_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRestaurantResponseMenuItem); i {
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
		file_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRestaurantResponseMenu); i {
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
			RawDescriptor: file_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
		MessageInfos:      file_service_proto_msgTypes,
	}.Build()
	File_service_proto = out.File
	file_service_proto_rawDesc = nil
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
