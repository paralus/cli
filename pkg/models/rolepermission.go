package models

type RolePermission struct {
	ApiVersion string             `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string             `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   Metadata           `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec       RolePermissionSpec `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty"`
	Status     Status             `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

type RolePermissionSpec struct {
	Permissions []string `protobuf:"bytes,1,rep,name=permissions,proto3" json:"permissions,omitempty"`
}

type RolePermissionList struct {
	ApiVersion string            `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string            `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   ListMetadata      `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Items      []*RolePermission `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
}
