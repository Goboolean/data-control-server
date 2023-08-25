// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.21.12
// source: api/redis-model/model.proto

package redis

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

type RedisStockAggregate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StockId   string  `protobuf:"bytes,1,opt,name=stockId,proto3" json:"stockId,omitempty"`
	EventType string  `protobuf:"bytes,2,opt,name=eventType,proto3" json:"eventType,omitempty"`
	Avg       float64 `protobuf:"fixed64,3,opt,name=avg,proto3" json:"avg,omitempty"`
	Min       float64 `protobuf:"fixed64,4,opt,name=min,proto3" json:"min,omitempty"`
	Max       float64 `protobuf:"fixed64,5,opt,name=max,proto3" json:"max,omitempty"`
	Start     float64 `protobuf:"fixed64,6,opt,name=start,proto3" json:"start,omitempty"`
	End       float64 `protobuf:"fixed64,7,opt,name=end,proto3" json:"end,omitempty"`
	StartTime int64   `protobuf:"varint,8,opt,name=startTime,proto3" json:"startTime,omitempty"`
	EndTime   int64   `protobuf:"varint,9,opt,name=endTime,proto3" json:"endTime,omitempty"`
}

func (x *RedisStockAggregate) Reset() {
	*x = RedisStockAggregate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_redis_model_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RedisStockAggregate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RedisStockAggregate) ProtoMessage() {}

func (x *RedisStockAggregate) ProtoReflect() protoreflect.Message {
	mi := &file_api_redis_model_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RedisStockAggregate.ProtoReflect.Descriptor instead.
func (*RedisStockAggregate) Descriptor() ([]byte, []int) {
	return file_api_redis_model_model_proto_rawDescGZIP(), []int{0}
}

func (x *RedisStockAggregate) GetStockId() string {
	if x != nil {
		return x.StockId
	}
	return ""
}

func (x *RedisStockAggregate) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

func (x *RedisStockAggregate) GetAvg() float64 {
	if x != nil {
		return x.Avg
	}
	return 0
}

func (x *RedisStockAggregate) GetMin() float64 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *RedisStockAggregate) GetMax() float64 {
	if x != nil {
		return x.Max
	}
	return 0
}

func (x *RedisStockAggregate) GetStart() float64 {
	if x != nil {
		return x.Start
	}
	return 0
}

func (x *RedisStockAggregate) GetEnd() float64 {
	if x != nil {
		return x.End
	}
	return 0
}

func (x *RedisStockAggregate) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *RedisStockAggregate) GetEndTime() int64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

var File_api_redis_model_model_proto protoreflect.FileDescriptor

var file_api_redis_model_model_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x65, 0x64, 0x69, 0x73, 0x2d, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x61,
	0x70, 0x69, 0x22, 0xe3, 0x01, 0x0a, 0x13, 0x52, 0x65, 0x64, 0x69, 0x73, 0x53, 0x74, 0x6f, 0x63,
	0x6b, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x74,
	0x6f, 0x63, 0x6b, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x74, 0x6f,
	0x63, 0x6b, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x76, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x03, 0x61, 0x76, 0x67, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x10,
	0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x65, 0x6e, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x41, 0x5a, 0x3f, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x47, 0x6f, 0x62, 0x6f, 0x6f, 0x6c, 0x65, 0x61, 0x6e,
	0x2f, 0x66, 0x65, 0x74, 0x63, 0x68, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x73, 0x74, 0x72, 0x75,
	0x63, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x72, 0x65, 0x64, 0x69, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_api_redis_model_model_proto_rawDescOnce sync.Once
	file_api_redis_model_model_proto_rawDescData = file_api_redis_model_model_proto_rawDesc
)

func file_api_redis_model_model_proto_rawDescGZIP() []byte {
	file_api_redis_model_model_proto_rawDescOnce.Do(func() {
		file_api_redis_model_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_redis_model_model_proto_rawDescData)
	})
	return file_api_redis_model_model_proto_rawDescData
}

var file_api_redis_model_model_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_api_redis_model_model_proto_goTypes = []interface{}{
	(*RedisStockAggregate)(nil), // 0: api.RedisStockAggregate
}
var file_api_redis_model_model_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_api_redis_model_model_proto_init() }
func file_api_redis_model_model_proto_init() {
	if File_api_redis_model_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_redis_model_model_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RedisStockAggregate); i {
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
			RawDescriptor: file_api_redis_model_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_redis_model_model_proto_goTypes,
		DependencyIndexes: file_api_redis_model_model_proto_depIdxs,
		MessageInfos:      file_api_redis_model_model_proto_msgTypes,
	}.Build()
	File_api_redis_model_model_proto = out.File
	file_api_redis_model_model_proto_rawDesc = nil
	file_api_redis_model_model_proto_goTypes = nil
	file_api_redis_model_model_proto_depIdxs = nil
}
