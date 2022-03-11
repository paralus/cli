package models

type Project struct {
	ApiVersion string      `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string      `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   Metadata    `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec       ProjectSpec `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty"`
	Status     Status      `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

type ProjectSpec struct {
	Default bool `protobuf:"varint,1,opt,name=default,proto3" json:"default,omitempty"`
}

type ProjectList struct {
	ApiVersion string       `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string       `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   ListMetadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Items      []*Project   `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
}
