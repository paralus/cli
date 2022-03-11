package models

type Role struct {
	ApiVersion string   `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string   `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   Metadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec       RoleSpec `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty"`
	Status     Status   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

type RoleSpec struct {
	Rolepermissions []string `protobuf:"bytes,1,rep,name=rolepermissions,proto3" json:"rolepermissions,omitempty"`
	IsGlobal        bool     `protobuf:"varint,2,opt,name=isGlobal,proto3" json:"isGlobal,omitempty"`
	Scope           string   `protobuf:"bytes,3,opt,name=scope,proto3" json:"scope,omitempty"`
}

type RoleList struct {
	ApiVersion string       `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string       `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   ListMetadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Items      []*Role      `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
}
