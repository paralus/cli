package models

type Idp struct {
	ApiVersion string   `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string   `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   Metadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec       IdpSpec  `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty"`
	Status     Status   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

type IdpSpec struct {
	IdpName            string `protobuf:"bytes,1,opt,name=idpName,proto3" json:"idpName,omitempty"`
	Domain             string `protobuf:"bytes,2,opt,name=domain,proto3" json:"domain,omitempty"`
	AcsUrl             string `protobuf:"bytes,3,opt,name=acsUrl,proto3" json:"acsUrl,omitempty"`
	SsoUrl             string `protobuf:"bytes,4,opt,name=ssoUrl,proto3" json:"ssoUrl,omitempty"`
	IdpCert            string `protobuf:"bytes,5,opt,name=idpCert,proto3" json:"idpCert,omitempty"`
	SpCert             string `protobuf:"bytes,6,opt,name=spCert,proto3" json:"spCert,omitempty"`
	MetadataUrl        string `protobuf:"bytes,7,opt,name=metadataUrl,proto3" json:"metadataUrl,omitempty"`
	MetadataFilename   string `protobuf:"bytes,8,opt,name=metadataFilename,proto3" json:"metadataFilename,omitempty"`
	SaeEnabled         bool   `protobuf:"varint,9,opt,name=saeEnabled,proto3" json:"saeEnabled,omitempty"`
	GroupAttributeName string `protobuf:"bytes,10,opt,name=groupAttributeName,proto3" json:"groupAttributeName,omitempty"`
	NameIdFormat       string `protobuf:"bytes,11,opt,name=nameIdFormat,proto3" json:"nameIdFormat,omitempty"`
	ConsumerBinding    string `protobuf:"bytes,12,opt,name=consumerBinding,proto3" json:"consumerBinding,omitempty"`
	SpEntityId         string `protobuf:"bytes,13,opt,name=spEntityId,proto3" json:"spEntityId,omitempty"`
}

type IdpList struct {
	ApiVersion string       `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string       `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   ListMetadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Items      []*Idp       `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
}
