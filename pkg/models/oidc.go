package models

import "google.golang.org/protobuf/types/known/structpb"

type OIDCProvider struct {
	ApiVersion string           `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string           `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   Metadata         `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec       OIDCProviderSpec `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty"`
	Status     Status           `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

type OIDCProviderSpec struct {
	ProviderName    string           `protobuf:"bytes,1,opt,name=providerName,proto3" json:"providerName,omitempty"` // enumeration?
	MapperUrl       string           `protobuf:"bytes,2,opt,name=mapperUrl,proto3" json:"mapperUrl,omitempty"`
	MapperFilename  string           `protobuf:"bytes,3,opt,name=mapperFilename,proto3" json:"mapperFilename,omitempty"`
	ClientId        string           `protobuf:"bytes,4,opt,name=clientId,proto3" json:"clientId,omitempty"`
	ClientSecret    string           `protobuf:"bytes,5,opt,name=clientSecret,proto3" json:"clientSecret,omitempty"`
	Scopes          []string         `protobuf:"bytes,6,rep,name=scopes,proto3" json:"scopes,omitempty"`
	IssuerUrl       string           `protobuf:"bytes,7,opt,name=issuerUrl,proto3" json:"issuerUrl,omitempty"`
	AuthUrl         string           `protobuf:"bytes,8,opt,name=authUrl,proto3" json:"authUrl,omitempty"`
	TokenUrl        string           `protobuf:"bytes,9,opt,name=tokenUrl,proto3" json:"tokenUrl,omitempty"`
	RequestedClaims *structpb.Struct `protobuf:"bytes,10,opt,name=requestedClaims,proto3" json:"requestedClaims,omitempty"` // JSON object
	Predefined      bool             `protobuf:"varint,11,opt,name=predefined,proto3" json:"predefined,omitempty"`
	CallbackUrl     string           `protobuf:"bytes,12,opt,name=callbackUrl,proto3" json:"callbackUrl,omitempty"`
}

type OIDCProviderList struct {
	ApiVersion string          `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string          `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   ListMetadata    `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Items      []*OIDCProvider `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
}
