// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: deviceReport.proto

package deviceReport

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// DeviceReportMessage 设备上报消息结构
type DeviceReportMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 消息id
	MsgId string `protobuf:"bytes,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	// 设备id
	Uuid string `protobuf:"bytes,2,opt,name=uuid,proto3" json:"uuid,omitempty"`
	// 产品码，paas为device_type
	ProductCode string `protobuf:"bytes,3,opt,name=product_code,json=productCode,proto3" json:"product_code,omitempty"`
	// 事件发生时间点
	EventTime int64 `protobuf:"varint,4,opt,name=event_time,json=eventTime,proto3" json:"event_time,omitempty"`
	// 数据内容
	Data *structpb.Struct `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	// 消息类型
	// DEVICE_REPORT 设备上报消息
	// SYSTEM_NOTIFY 系统消息
	Type string `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
	// 业务appId
	AppId int64 `protobuf:"varint,8,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"` // @gotags: validate:"required"
	// 消息上报次数，兼容旧版本
	Times string `protobuf:"bytes,9,opt,name=times,proto3" json:"times,omitempty"`
	// carrier
	Carrier []byte `protobuf:"bytes,10,opt,name=carrier,proto3" json:"carrier,omitempty"`
}

func (x *DeviceReportMessage) Reset() {
	*x = DeviceReportMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_deviceReport_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceReportMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceReportMessage) ProtoMessage() {}

func (x *DeviceReportMessage) ProtoReflect() protoreflect.Message {
	mi := &file_deviceReport_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceReportMessage.ProtoReflect.Descriptor instead.
func (*DeviceReportMessage) Descriptor() ([]byte, []int) {
	return file_deviceReport_proto_rawDescGZIP(), []int{0}
}

func (x *DeviceReportMessage) GetMsgId() string {
	if x != nil {
		return x.MsgId
	}
	return ""
}

func (x *DeviceReportMessage) GetUuid() string {
	if x != nil {
		return x.Uuid
	}
	return ""
}

func (x *DeviceReportMessage) GetProductCode() string {
	if x != nil {
		return x.ProductCode
	}
	return ""
}

func (x *DeviceReportMessage) GetEventTime() int64 {
	if x != nil {
		return x.EventTime
	}
	return 0
}

func (x *DeviceReportMessage) GetData() *structpb.Struct {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DeviceReportMessage) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *DeviceReportMessage) GetAppId() int64 {
	if x != nil {
		return x.AppId
	}
	return 0
}

func (x *DeviceReportMessage) GetTimes() string {
	if x != nil {
		return x.Times
	}
	return ""
}

func (x *DeviceReportMessage) GetCarrier() []byte {
	if x != nil {
		return x.Carrier
	}
	return nil
}

var File_deviceReport_proto protoreflect.FileDescriptor

var file_deviceReport_proto_rawDesc = []byte{
	0x0a, 0x12, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x8a, 0x02, 0x0a, 0x13, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x73,
	0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x75, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x75, 0x75, 0x69, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x5f,
	0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x63, 0x61, 0x72, 0x72, 0x69, 0x65, 0x72, 0x42,
	0x11, 0x5a, 0x0f, 0x2e, 0x2f, 0x3b, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_deviceReport_proto_rawDescOnce sync.Once
	file_deviceReport_proto_rawDescData = file_deviceReport_proto_rawDesc
)

func file_deviceReport_proto_rawDescGZIP() []byte {
	file_deviceReport_proto_rawDescOnce.Do(func() {
		file_deviceReport_proto_rawDescData = protoimpl.X.CompressGZIP(file_deviceReport_proto_rawDescData)
	})
	return file_deviceReport_proto_rawDescData
}

var file_deviceReport_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_deviceReport_proto_goTypes = []interface{}{
	(*DeviceReportMessage)(nil), // 0: DeviceReportMessage
	(*structpb.Struct)(nil),     // 1: google.protobuf.Struct
}
var file_deviceReport_proto_depIdxs = []int32{
	1, // 0: DeviceReportMessage.data:type_name -> google.protobuf.Struct
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_deviceReport_proto_init() }
func file_deviceReport_proto_init() {
	if File_deviceReport_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_deviceReport_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceReportMessage); i {
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
			RawDescriptor: file_deviceReport_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_deviceReport_proto_goTypes,
		DependencyIndexes: file_deviceReport_proto_depIdxs,
		MessageInfos:      file_deviceReport_proto_msgTypes,
	}.Build()
	File_deviceReport_proto = out.File
	file_deviceReport_proto_rawDesc = nil
	file_deviceReport_proto_goTypes = nil
	file_deviceReport_proto_depIdxs = nil
}