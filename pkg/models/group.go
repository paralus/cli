package models

type Group struct {
	ApiVersion string    `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string    `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   Metadata  `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec       GroupSpec `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty"`
	Status     Status    `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

type GroupSpec struct {
	ProjectNamespaceRoles []*ProjectNamespaceRole `protobuf:"bytes,1,rep,name=projectNamespaceRoles,proto3" json:"projectNamespaceRoles,omitempty"`
	Users                 []string                `protobuf:"bytes,2,rep,name=users,proto3" json:"users,omitempty"`
	Type                  string                  `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
}

type GroupList struct {
	ApiVersion string       `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string       `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   ListMetadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Items      []*Group     `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
}
