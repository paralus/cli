package models

type Metro struct {
	Id          string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	City        string `protobuf:"bytes,3,opt,name=city,proto3" json:"city,omitempty"`
	State       string `protobuf:"bytes,4,opt,name=state,proto3" json:"state,omitempty"`
	Country     string `protobuf:"bytes,5,opt,name=country,proto3" json:"country,omitempty"`
	Locale      string `protobuf:"bytes,6,opt,name=locale,proto3" json:"locale,omitempty"`
	Latitude    string `protobuf:"bytes,7,opt,name=latitude,proto3" json:"latitude,omitempty"`
	Longitude   string `protobuf:"bytes,8,opt,name=longitude,proto3" json:"longitude,omitempty"`
	CountryCode string `protobuf:"bytes,9,opt,name=countryCode,proto3" json:"countryCode,omitempty"`
	StateCode   string `protobuf:"bytes,10,opt,name=stateCode,proto3" json:"stateCode,omitempty"`
}

type Location struct {
	ApiVersion string   `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string   `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   Metadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Spec       Metro    `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty"`
	Status     Status   `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
}

type LocationList struct {
	ApiVersion string       `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty"`
	Kind       string       `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty"`
	Metadata   ListMetadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Items      []*Metro     `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty"`
}
