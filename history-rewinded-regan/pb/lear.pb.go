// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: lear.proto

package pb

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

type IncidentType int32

const (
	IncidentType_INCIDENT_TYPE_UNSPECIFIED IncidentType = 0
	IncidentType_INCIDENT_TYPE_EVENT       IncidentType = 1
	IncidentType_INCIDENT_TYPE_BIRTH       IncidentType = 2
	IncidentType_INCIDENT_TYPE_DEATH       IncidentType = 3
	IncidentType_INCIDENT_TYPE_HOLIDAY     IncidentType = 4
)

// Enum value maps for IncidentType.
var (
	IncidentType_name = map[int32]string{
		0: "INCIDENT_TYPE_UNSPECIFIED",
		1: "INCIDENT_TYPE_EVENT",
		2: "INCIDENT_TYPE_BIRTH",
		3: "INCIDENT_TYPE_DEATH",
		4: "INCIDENT_TYPE_HOLIDAY",
	}
	IncidentType_value = map[string]int32{
		"INCIDENT_TYPE_UNSPECIFIED": 0,
		"INCIDENT_TYPE_EVENT":       1,
		"INCIDENT_TYPE_BIRTH":       2,
		"INCIDENT_TYPE_DEATH":       3,
		"INCIDENT_TYPE_HOLIDAY":     4,
	}
)

func (x IncidentType) Enum() *IncidentType {
	p := new(IncidentType)
	*p = x
	return p
}

func (x IncidentType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (IncidentType) Descriptor() protoreflect.EnumDescriptor {
	return file_lear_proto_enumTypes[0].Descriptor()
}

func (IncidentType) Type() protoreflect.EnumType {
	return &file_lear_proto_enumTypes[0]
}

func (x IncidentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use IncidentType.Descriptor instead.
func (IncidentType) EnumDescriptor() ([]byte, []int) {
	return file_lear_proto_rawDescGZIP(), []int{0}
}

type FetchIncidentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Month int64 `protobuf:"varint,2,opt,name=month,proto3" json:"month,omitempty"`
	Day   int64 `protobuf:"varint,3,opt,name=day,proto3" json:"day,omitempty"`
}

func (x *FetchIncidentRequest) Reset() {
	*x = FetchIncidentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lear_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FetchIncidentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FetchIncidentRequest) ProtoMessage() {}

func (x *FetchIncidentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_lear_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FetchIncidentRequest.ProtoReflect.Descriptor instead.
func (*FetchIncidentRequest) Descriptor() ([]byte, []int) {
	return file_lear_proto_rawDescGZIP(), []int{0}
}

func (x *FetchIncidentRequest) GetMonth() int64 {
	if x != nil {
		return x.Month
	}
	return 0
}

func (x *FetchIncidentRequest) GetDay() int64 {
	if x != nil {
		return x.Day
	}
	return 0
}

type Incident struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Year             int64        `protobuf:"varint,1,opt,name=year,proto3" json:"year,omitempty"`
	Month            int64        `protobuf:"varint,2,opt,name=month,proto3" json:"month,omitempty"`
	Day              int64        `protobuf:"varint,3,opt,name=day,proto3" json:"day,omitempty"`
	IncidentType     IncidentType `protobuf:"varint,4,opt,name=incident_type,json=incidentType,proto3,enum=IncidentType" json:"incident_type,omitempty"`
	Summary          string       `protobuf:"bytes,5,opt,name=summary,proto3" json:"summary,omitempty"`
	IncidentInDetail string       `protobuf:"bytes,6,opt,name=incident_in_detail,json=incidentInDetail,proto3" json:"incident_in_detail,omitempty"`
}

func (x *Incident) Reset() {
	*x = Incident{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lear_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Incident) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Incident) ProtoMessage() {}

func (x *Incident) ProtoReflect() protoreflect.Message {
	mi := &file_lear_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Incident.ProtoReflect.Descriptor instead.
func (*Incident) Descriptor() ([]byte, []int) {
	return file_lear_proto_rawDescGZIP(), []int{1}
}

func (x *Incident) GetYear() int64 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *Incident) GetMonth() int64 {
	if x != nil {
		return x.Month
	}
	return 0
}

func (x *Incident) GetDay() int64 {
	if x != nil {
		return x.Day
	}
	return 0
}

func (x *Incident) GetIncidentType() IncidentType {
	if x != nil {
		return x.IncidentType
	}
	return IncidentType_INCIDENT_TYPE_UNSPECIFIED
}

func (x *Incident) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *Incident) GetIncidentInDetail() string {
	if x != nil {
		return x.IncidentInDetail
	}
	return ""
}

var File_lear_proto protoreflect.FileDescriptor

var file_lear_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6c, 0x65, 0x61, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x14,
	0x46, 0x65, 0x74, 0x63, 0x68, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x61,
	0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x64, 0x61, 0x79, 0x22, 0xc2, 0x01, 0x0a,
	0x08, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x79, 0x65, 0x61,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x79, 0x65, 0x61, 0x72, 0x12, 0x14, 0x0a,
	0x05, 0x6d, 0x6f, 0x6e, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6d, 0x6f,
	0x6e, 0x74, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x64, 0x61, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x03, 0x64, 0x61, 0x79, 0x12, 0x32, 0x0a, 0x0d, 0x69, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0d, 0x2e, 0x49,
	0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x69, 0x6e, 0x63,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d,
	0x61, 0x72, 0x79, 0x12, 0x2c, 0x0a, 0x12, 0x69, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x5f,
	0x69, 0x6e, 0x5f, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x10, 0x69, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x44, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x2a, 0x93, 0x01, 0x0a, 0x0c, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1d, 0x0a, 0x19, 0x49, 0x4e, 0x43, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x17, 0x0a, 0x13, 0x49, 0x4e, 0x43, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x49, 0x4e,
	0x43, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x49, 0x52, 0x54,
	0x48, 0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x49, 0x4e, 0x43, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x5f,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x45, 0x41, 0x54, 0x48, 0x10, 0x03, 0x12, 0x19, 0x0a, 0x15,
	0x49, 0x4e, 0x43, 0x49, 0x44, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x48, 0x4f,
	0x4c, 0x49, 0x44, 0x41, 0x59, 0x10, 0x04, 0x32, 0x94, 0x02, 0x0a, 0x04, 0x4c, 0x65, 0x61, 0x72,
	0x12, 0x36, 0x0a, 0x10, 0x46, 0x65, 0x74, 0x63, 0x68, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e,
	0x74, 0x73, 0x4f, 0x6e, 0x12, 0x15, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x49, 0x6e, 0x63, 0x69,
	0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x49, 0x6e,
	0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x30, 0x01, 0x12, 0x33, 0x0a, 0x0d, 0x46, 0x65, 0x74, 0x63,
	0x68, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x4f, 0x6e, 0x12, 0x15, 0x2e, 0x46, 0x65, 0x74, 0x63,
	0x68, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x09, 0x2e, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x30, 0x01, 0x12, 0x33, 0x0a,
	0x0d, 0x46, 0x65, 0x74, 0x63, 0x68, 0x42, 0x69, 0x72, 0x74, 0x68, 0x73, 0x4f, 0x6e, 0x12, 0x15,
	0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x30, 0x01, 0x12, 0x33, 0x0a, 0x0d, 0x46, 0x65, 0x74, 0x63, 0x68, 0x44, 0x65, 0x61, 0x74, 0x68,
	0x73, 0x4f, 0x6e, 0x12, 0x15, 0x2e, 0x46, 0x65, 0x74, 0x63, 0x68, 0x49, 0x6e, 0x63, 0x69, 0x64,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x49, 0x6e, 0x63,
	0x69, 0x64, 0x65, 0x6e, 0x74, 0x30, 0x01, 0x12, 0x35, 0x0a, 0x0f, 0x46, 0x65, 0x74, 0x63, 0x68,
	0x48, 0x6f, 0x6c, 0x69, 0x64, 0x61, 0x79, 0x73, 0x4f, 0x6e, 0x12, 0x15, 0x2e, 0x46, 0x65, 0x74,
	0x63, 0x68, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x09, 0x2e, 0x49, 0x6e, 0x63, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x30, 0x01, 0x42, 0x07,
	0x5a, 0x05, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lear_proto_rawDescOnce sync.Once
	file_lear_proto_rawDescData = file_lear_proto_rawDesc
)

func file_lear_proto_rawDescGZIP() []byte {
	file_lear_proto_rawDescOnce.Do(func() {
		file_lear_proto_rawDescData = protoimpl.X.CompressGZIP(file_lear_proto_rawDescData)
	})
	return file_lear_proto_rawDescData
}

var file_lear_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_lear_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_lear_proto_goTypes = []interface{}{
	(IncidentType)(0),            // 0: IncidentType
	(*FetchIncidentRequest)(nil), // 1: FetchIncidentRequest
	(*Incident)(nil),             // 2: Incident
}
var file_lear_proto_depIdxs = []int32{
	0, // 0: Incident.incident_type:type_name -> IncidentType
	1, // 1: Lear.FetchIncidentsOn:input_type -> FetchIncidentRequest
	1, // 2: Lear.FetchEventsOn:input_type -> FetchIncidentRequest
	1, // 3: Lear.FetchBirthsOn:input_type -> FetchIncidentRequest
	1, // 4: Lear.FetchDeathsOn:input_type -> FetchIncidentRequest
	1, // 5: Lear.FetchHolidaysOn:input_type -> FetchIncidentRequest
	2, // 6: Lear.FetchIncidentsOn:output_type -> Incident
	2, // 7: Lear.FetchEventsOn:output_type -> Incident
	2, // 8: Lear.FetchBirthsOn:output_type -> Incident
	2, // 9: Lear.FetchDeathsOn:output_type -> Incident
	2, // 10: Lear.FetchHolidaysOn:output_type -> Incident
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_lear_proto_init() }
func file_lear_proto_init() {
	if File_lear_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lear_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FetchIncidentRequest); i {
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
		file_lear_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Incident); i {
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
			RawDescriptor: file_lear_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_lear_proto_goTypes,
		DependencyIndexes: file_lear_proto_depIdxs,
		EnumInfos:         file_lear_proto_enumTypes,
		MessageInfos:      file_lear_proto_msgTypes,
	}.Build()
	File_lear_proto = out.File
	file_lear_proto_rawDesc = nil
	file_lear_proto_goTypes = nil
	file_lear_proto_depIdxs = nil
}
