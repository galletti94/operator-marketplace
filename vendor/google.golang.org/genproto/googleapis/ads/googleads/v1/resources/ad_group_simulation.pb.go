// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v1/resources/ad_group_simulation.proto

package resources

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	common "google.golang.org/genproto/googleapis/ads/googleads/v1/common"
	enums "google.golang.org/genproto/googleapis/ads/googleads/v1/enums"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// An ad group simulation. Supported combinations of advertising
// channel type, simulation type and simulation modification method is
// detailed below respectively.
//
// SEARCH   CPC_BID     DEFAULT
// SEARCH   CPC_BID     UNIFORM
// SEARCH   TARGET_CPA  UNIFORM
// DISPLAY  CPC_BID     DEFAULT
// DISPLAY  CPC_BID     UNIFORM
// DISPLAY  TARGET_CPA  UNIFORM
type AdGroupSimulation struct {
	// The resource name of the ad group simulation.
	// Ad group simulation resource names have the form:
	//
	//
	// `customers/{customer_id}/adGroupSimulations/{ad_group_id}~{type}~{modification_method}~{start_date}~{end_date}`
	ResourceName string `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	// Ad group id of the simulation.
	AdGroupId *wrappers.Int64Value `protobuf:"bytes,2,opt,name=ad_group_id,json=adGroupId,proto3" json:"ad_group_id,omitempty"`
	// The field that the simulation modifies.
	Type enums.SimulationTypeEnum_SimulationType `protobuf:"varint,3,opt,name=type,proto3,enum=google.ads.googleads.v1.enums.SimulationTypeEnum_SimulationType" json:"type,omitempty"`
	// How the simulation modifies the field.
	ModificationMethod enums.SimulationModificationMethodEnum_SimulationModificationMethod `protobuf:"varint,4,opt,name=modification_method,json=modificationMethod,proto3,enum=google.ads.googleads.v1.enums.SimulationModificationMethodEnum_SimulationModificationMethod" json:"modification_method,omitempty"`
	// First day on which the simulation is based, in YYYY-MM-DD format.
	StartDate *wrappers.StringValue `protobuf:"bytes,5,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	// Last day on which the simulation is based, in YYYY-MM-DD format
	EndDate *wrappers.StringValue `protobuf:"bytes,6,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	// List of simulation points.
	//
	// Types that are valid to be assigned to PointList:
	//	*AdGroupSimulation_CpcBidPointList
	//	*AdGroupSimulation_TargetCpaPointList
	PointList            isAdGroupSimulation_PointList `protobuf_oneof:"point_list"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *AdGroupSimulation) Reset()         { *m = AdGroupSimulation{} }
func (m *AdGroupSimulation) String() string { return proto.CompactTextString(m) }
func (*AdGroupSimulation) ProtoMessage()    {}
func (*AdGroupSimulation) Descriptor() ([]byte, []int) {
	return fileDescriptor_8c8573d5f5d3a023, []int{0}
}

func (m *AdGroupSimulation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AdGroupSimulation.Unmarshal(m, b)
}
func (m *AdGroupSimulation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AdGroupSimulation.Marshal(b, m, deterministic)
}
func (m *AdGroupSimulation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AdGroupSimulation.Merge(m, src)
}
func (m *AdGroupSimulation) XXX_Size() int {
	return xxx_messageInfo_AdGroupSimulation.Size(m)
}
func (m *AdGroupSimulation) XXX_DiscardUnknown() {
	xxx_messageInfo_AdGroupSimulation.DiscardUnknown(m)
}

var xxx_messageInfo_AdGroupSimulation proto.InternalMessageInfo

func (m *AdGroupSimulation) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func (m *AdGroupSimulation) GetAdGroupId() *wrappers.Int64Value {
	if m != nil {
		return m.AdGroupId
	}
	return nil
}

func (m *AdGroupSimulation) GetType() enums.SimulationTypeEnum_SimulationType {
	if m != nil {
		return m.Type
	}
	return enums.SimulationTypeEnum_UNSPECIFIED
}

func (m *AdGroupSimulation) GetModificationMethod() enums.SimulationModificationMethodEnum_SimulationModificationMethod {
	if m != nil {
		return m.ModificationMethod
	}
	return enums.SimulationModificationMethodEnum_UNSPECIFIED
}

func (m *AdGroupSimulation) GetStartDate() *wrappers.StringValue {
	if m != nil {
		return m.StartDate
	}
	return nil
}

func (m *AdGroupSimulation) GetEndDate() *wrappers.StringValue {
	if m != nil {
		return m.EndDate
	}
	return nil
}

type isAdGroupSimulation_PointList interface {
	isAdGroupSimulation_PointList()
}

type AdGroupSimulation_CpcBidPointList struct {
	CpcBidPointList *common.CpcBidSimulationPointList `protobuf:"bytes,8,opt,name=cpc_bid_point_list,json=cpcBidPointList,proto3,oneof"`
}

type AdGroupSimulation_TargetCpaPointList struct {
	TargetCpaPointList *common.TargetCpaSimulationPointList `protobuf:"bytes,9,opt,name=target_cpa_point_list,json=targetCpaPointList,proto3,oneof"`
}

func (*AdGroupSimulation_CpcBidPointList) isAdGroupSimulation_PointList() {}

func (*AdGroupSimulation_TargetCpaPointList) isAdGroupSimulation_PointList() {}

func (m *AdGroupSimulation) GetPointList() isAdGroupSimulation_PointList {
	if m != nil {
		return m.PointList
	}
	return nil
}

func (m *AdGroupSimulation) GetCpcBidPointList() *common.CpcBidSimulationPointList {
	if x, ok := m.GetPointList().(*AdGroupSimulation_CpcBidPointList); ok {
		return x.CpcBidPointList
	}
	return nil
}

func (m *AdGroupSimulation) GetTargetCpaPointList() *common.TargetCpaSimulationPointList {
	if x, ok := m.GetPointList().(*AdGroupSimulation_TargetCpaPointList); ok {
		return x.TargetCpaPointList
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*AdGroupSimulation) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*AdGroupSimulation_CpcBidPointList)(nil),
		(*AdGroupSimulation_TargetCpaPointList)(nil),
	}
}

func init() {
	proto.RegisterType((*AdGroupSimulation)(nil), "google.ads.googleads.v1.resources.AdGroupSimulation")
}

func init() {
	proto.RegisterFile("google/ads/googleads/v1/resources/ad_group_simulation.proto", fileDescriptor_8c8573d5f5d3a023)
}

var fileDescriptor_8c8573d5f5d3a023 = []byte{
	// 563 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x54, 0xdd, 0x6e, 0xd3, 0x3c,
	0x18, 0xfe, 0xd2, 0xfd, 0x7c, 0x9b, 0x37, 0x40, 0x18, 0x81, 0xa2, 0x31, 0xa1, 0x0e, 0x34, 0xa9,
	0x47, 0x8e, 0xba, 0x21, 0x10, 0x29, 0x07, 0xa4, 0x03, 0x95, 0x21, 0x86, 0xaa, 0xac, 0xea, 0x01,
	0xaa, 0x14, 0xb9, 0xb1, 0x97, 0x59, 0x6a, 0x6c, 0x13, 0x3b, 0x43, 0x3b, 0xe6, 0x98, 0x9b, 0xe0,
	0x90, 0x6b, 0xe0, 0x0a, 0xb8, 0x14, 0xae, 0x02, 0xc5, 0x49, 0xdc, 0xaa, 0x5d, 0x59, 0xcf, 0x5e,
	0xbf, 0xef, 0xf3, 0x63, 0x3f, 0xb1, 0x03, 0x3a, 0x89, 0x10, 0xc9, 0x84, 0x7a, 0x98, 0x28, 0xaf,
	0x2c, 0x8b, 0xea, 0xaa, 0xed, 0x65, 0x54, 0x89, 0x3c, 0x8b, 0xa9, 0xf2, 0x30, 0x89, 0x92, 0x4c,
	0xe4, 0x32, 0x52, 0x2c, 0xcd, 0x27, 0x58, 0x33, 0xc1, 0x91, 0xcc, 0x84, 0x16, 0xf0, 0xa0, 0x64,
	0x20, 0x4c, 0x14, 0xb2, 0x64, 0x74, 0xd5, 0x46, 0x96, 0xbc, 0xe7, 0x2d, 0xd3, 0x8f, 0x45, 0x9a,
	0x0a, 0xee, 0xcd, 0x6b, 0xee, 0x75, 0x97, 0x11, 0x28, 0xcf, 0x53, 0x35, 0x83, 0x8f, 0x52, 0x41,
	0xd8, 0x05, 0x8b, 0xab, 0x05, 0xd5, 0x97, 0x82, 0x54, 0x1a, 0xc7, 0x2b, 0x6b, 0xe8, 0x6b, 0x49,
	0x2b, 0xd2, 0x93, 0x8a, 0x64, 0x56, 0xe3, 0xfc, 0xc2, 0xfb, 0x9a, 0x61, 0x29, 0x69, 0xa6, 0xaa,
	0xf9, 0x7e, 0x2d, 0x2a, 0x99, 0x87, 0x39, 0x17, 0xda, 0x28, 0x54, 0xd3, 0xa7, 0xbf, 0x36, 0xc0,
	0xfd, 0x80, 0xf4, 0x8a, 0x9c, 0xce, 0xad, 0x3c, 0x7c, 0x06, 0xee, 0xd4, 0x51, 0x44, 0x1c, 0xa7,
	0xd4, 0x75, 0x9a, 0x4e, 0x6b, 0x3b, 0xdc, 0xad, 0x9b, 0x9f, 0x70, 0x4a, 0x61, 0x07, 0xec, 0xd8,
	0x88, 0x19, 0x71, 0x1b, 0x4d, 0xa7, 0xb5, 0x73, 0xf4, 0xb8, 0x0a, 0x14, 0xd5, 0xdb, 0x41, 0xa7,
	0x5c, 0xbf, 0x78, 0x3e, 0xc4, 0x93, 0x9c, 0x86, 0xdb, 0xb8, 0x74, 0x3a, 0x25, 0x70, 0x00, 0xd6,
	0x8b, 0x33, 0xb8, 0x6b, 0x4d, 0xa7, 0x75, 0xf7, 0xe8, 0x0d, 0x5a, 0xf6, 0x45, 0xcc, 0xc9, 0xd1,
	0x74, 0x6b, 0x83, 0x6b, 0x49, 0xdf, 0xf1, 0x3c, 0x9d, 0x6b, 0x85, 0x46, 0x0d, 0x7e, 0x77, 0xc0,
	0x83, 0x1b, 0xe2, 0x75, 0xd7, 0x8d, 0xcb, 0x68, 0x65, 0x97, 0xb3, 0x19, 0x8d, 0x33, 0x23, 0x31,
	0xe7, 0xb9, 0x08, 0x08, 0x61, 0xba, 0xd0, 0x83, 0x1d, 0x00, 0x94, 0xc6, 0x99, 0x8e, 0x08, 0xd6,
	0xd4, 0xdd, 0x30, 0x09, 0xed, 0x2f, 0x24, 0x74, 0xae, 0x33, 0xc6, 0x93, 0x2a, 0x22, 0x83, 0x7f,
	0x8b, 0x35, 0x85, 0x2f, 0xc1, 0x16, 0xe5, 0xa4, 0xa4, 0x6e, 0xae, 0x40, 0xfd, 0x9f, 0x72, 0x62,
	0x88, 0x97, 0x00, 0xc6, 0x32, 0x8e, 0xc6, 0x8c, 0x44, 0x52, 0x30, 0xae, 0xa3, 0x09, 0x53, 0xda,
	0xdd, 0x32, 0x12, 0xaf, 0x96, 0x66, 0x50, 0x5e, 0x6c, 0x74, 0x22, 0xe3, 0x2e, 0x23, 0xd3, 0x93,
	0xf6, 0x0b, 0x85, 0x8f, 0x4c, 0xe9, 0xf7, 0xff, 0x85, 0xf7, 0x62, 0x33, 0xb4, 0x2d, 0xf8, 0x05,
	0x3c, 0xd4, 0x38, 0x4b, 0xa8, 0x8e, 0x62, 0x89, 0x67, 0xcd, 0xb6, 0x8d, 0xd9, 0xeb, 0xdb, 0xcc,
	0x06, 0x86, 0x7c, 0x22, 0xf1, 0xcd, 0x7e, 0x50, 0xd7, 0x73, 0xdb, 0xed, 0xee, 0x02, 0x30, 0xf5,
	0xe9, 0x7e, 0x6b, 0x80, 0xc3, 0x58, 0xa4, 0xe8, 0xd6, 0x07, 0xdd, 0x7d, 0xb4, 0x70, 0xcb, 0xfb,
	0x45, 0x88, 0x7d, 0xe7, 0xf3, 0x87, 0x8a, 0x9c, 0x88, 0x09, 0xe6, 0x09, 0x12, 0x59, 0xe2, 0x25,
	0x94, 0x9b, 0x88, 0xeb, 0x57, 0x28, 0x99, 0xfa, 0xc7, 0x9f, 0xa6, 0x63, 0xab, 0x1f, 0x8d, 0xb5,
	0x5e, 0x10, 0xfc, 0x6c, 0x1c, 0xf4, 0x4a, 0xc9, 0x80, 0x28, 0x54, 0x96, 0x45, 0x35, 0x6c, 0xa3,
	0xb0, 0x46, 0xfe, 0xae, 0x31, 0xa3, 0x80, 0xa8, 0x91, 0xc5, 0x8c, 0x86, 0xed, 0x91, 0xc5, 0xfc,
	0x69, 0x1c, 0x96, 0x03, 0xdf, 0x0f, 0x88, 0xf2, 0x7d, 0x8b, 0xf2, 0xfd, 0x61, 0xdb, 0xf7, 0x2d,
	0x6e, 0xbc, 0x69, 0x36, 0x7b, 0xfc, 0x37, 0x00, 0x00, 0xff, 0xff, 0x05, 0xed, 0xb9, 0x15, 0x15,
	0x05, 0x00, 0x00,
}
