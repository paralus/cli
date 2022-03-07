package models

type User struct {
	ApiVersion string   `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string   `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   Metadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec       UserSpec `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty"`
	Status     Status   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

type UserSpec struct {
	FirstName             string                  `protobuf:"bytes,1,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName              string                  `protobuf:"bytes,2,opt,name=lastName,proto3" json:"lastName,omitempty"`
	Phone                 string                  `protobuf:"bytes,4,opt,name=phone,proto3" json:"phone,omitempty"`
	Password              string                  `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	Groups                []string                `protobuf:"bytes,6,rep,name=groups,proto3" json:"groups,omitempty"`
	ProjectNamespaceRoles []*ProjectNamespaceRole `protobuf:"bytes,7,rep,name=projectNamespaceRoles,proto3" json:"projectNamespaceRoles,omitempty"`
	EmailVerified         bool                    `protobuf:"varint,8,opt,name=emailVerified,proto3" json:"emailVerified,omitempty"`
	PhoneVerified         bool                    `protobuf:"varint,9,opt,name=phoneVerified,proto3" json:"phoneVerified,omitempty"`
}

type ProjectNamespaceRole struct {
	Project   string `protobuf:"bytes,1,opt,name=project,proto3,oneof" json:"project,omitempty"`
	Namespace int64  `protobuf:"varint,2,opt,name=namespace,proto3,oneof" json:"namespace,omitempty"`
	Role      string `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
}

type UserList struct {
	ApiVersion string       `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string       `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   ListMetadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Items      []*User      `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
}
