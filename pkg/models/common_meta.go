package models

type Metadata struct {
	Name         string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	DisplayName  string            `protobuf:"bytes,2,opt,name=displayName,proto3" json:"displayName,omitempty"`
	Description  string            `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Labels       map[string]string `protobuf:"bytes,4,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Annotations  map[string]string `protobuf:"bytes,5,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Project      string            `protobuf:"bytes,6,opt,name=project,proto3" json:"project,omitempty"`
	Organization string            `protobuf:"bytes,7,opt,name=organization,proto3" json:"organization,omitempty"`
	Partner      string            `protobuf:"bytes,8,opt,name=partner,proto3" json:"partner,omitempty"`
	Id           string            `protobuf:"bytes,9,opt,name=id,proto3" json:"id,omitempty"`
}

type Status struct {
	ConditionType   string          `protobuf:"bytes,1,opt,name=conditionType,proto3" json:"conditionType,omitempty"`
	ConditionStatus ConditionStatus `protobuf:"varint,2,opt,name=conditionStatus,proto3,enum=rafay.dev.types.common.v3.ConditionStatus" json:"conditionStatus,omitempty"`
	Reason          string          `protobuf:"bytes,4,opt,name=reason,proto3" json:"reason,omitempty"`
}

type ConditionStatus int32

const (
	ConditionStatus_StatusNotSet    ConditionStatus = 0
	ConditionStatus_StatusSubmitted ConditionStatus = 1
	ConditionStatus_StatusOK        ConditionStatus = 2
	ConditionStatus_StatusFailed    ConditionStatus = 3
)

// Enum value maps for ConditionStatus.
var (
	ConditionStatus_name = map[int32]string{
		0: "StatusNotSet",
		1: "StatusSubmitted",
		2: "StatusOK",
		3: "StatusFailed",
	}
	ConditionStatus_value = map[string]int32{
		"StatusNotSet":    0,
		"StatusSubmitted": 1,
		"StatusOK":        2,
		"StatusFailed":    3,
	}
)

func (x ConditionStatus) Enum() *ConditionStatus {
	p := new(ConditionStatus)
	*p = x
	return p
}

type ListMetadata struct {
	Count  int64 `protobuf:"zigzag64,1,opt,name=count,proto3" json:"count,omitempty"`
	Offset int64 `protobuf:"zigzag64,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Limit  int64 `protobuf:"zigzag64,3,opt,name=limit,proto3" json:"limit,omitempty"`
}
