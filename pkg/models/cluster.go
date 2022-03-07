package models

import (
	"time"
)

const (
	ClusterTypeImport = "imported"
)

type Cluster struct {
	ApiVersion string      `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string      `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty" yaml:"kind"`
	Metadata   Metadata    `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty" yaml:"metadata"`
	Spec       ClusterSpec `protobuf:"bytes,4,opt,name=spec,proto3" json:"spec,omitempty" yaml:"spec"`
	Status     Status      `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty" yaml:"status"`
}

type ClusterSpec struct {
	ClusterType      string           `protobuf:"bytes,1,opt,name=clusterType,proto3" json:"clusterType,omitempty" yaml:"clusterType"`
	Metro            *Metro           `protobuf:"bytes,2,opt,name=metro,proto3" json:"metro,omitempty" yaml:"metro"`
	OverrideSelector string           `protobuf:"bytes,3,opt,name=overrideSelector,proto3" json:"overrideSelector,omitempty" yaml:"overrideSelector"`
	Params           *ProvisionParams `protobuf:"bytes,4,opt,name=params,proto3" json:"params,omitempty" yaml:"params"`
	ShareMode        ClusterShareMode `protobuf:"varint,5,opt,name=shareMode,proto3,enum=rafay.dev.types.infra.v3.ClusterShareMode" json:"shareMode,omitempty" yaml:"shareMode"`
	ProxyConfig      *ProxyConfig     `protobuf:"bytes,6,opt,name=proxyConfig,proto3" json:"proxyConfig,omitempty" yaml:"proxyConfig"`
	ClusterData      *ClusterData     `protobuf:"bytes,7,opt,name=clusterData,proto3" json:"clusterData,omitempty" yaml:"clusterData"`
}

type ClusterData struct {
	Provider      string            `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty" yaml:"provider"`
	Passphrase    string            `protobuf:"bytes,2,opt,name=passphrase,proto3" json:"passphrase,omitempty" yaml:"passphrase"`
	Cname         string            `protobuf:"bytes,3,opt,name=cname,proto3" json:"cname,omitempty" yaml:"cname"`
	Arecord       string            `protobuf:"bytes,4,opt,name=arecord,proto3" json:"arecord,omitempty" yaml:"arecord"`
	DisplayName   string            `protobuf:"bytes,5,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty" yaml:"displayName"`
	Health        Health            `protobuf:"varint,6,opt,name=health,proto3,enum=rafay.dev.types.infra.v3.Health" json:"health,omitempty" yaml:"health"`
	Manufacturer  string            `protobuf:"bytes,7,opt,name=manufacturer,proto3" json:"manufacturer,omitempty" yaml:"manufacturer"`
	Projects      []*ProjectCluster `protobuf:"bytes,10,rep,name=projects,proto3" json:"projects,omitempty" yaml:"projects"`
	ClusterStatus *ClusterStatus    `protobuf:"bytes,11,opt,name=cluster_status,json=clusterStatus,proto3" json:"cluster_status,omitempty" yaml:"clusterStatus"`
}

type ClusterStatus struct {
	Conditions         []*ClusterCondition `protobuf:"bytes,1,rep,name=conditions,proto3" json:"conditions,omitempty" yaml:"conditions"`
	Token              string              `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty" yaml:"token"`
	SystemTaskCount    int64               `protobuf:"zigzag64,4,opt,name=systemTaskCount,proto3" json:"systemTaskCount,omitempty" yaml:"systemTaskCount"`
	CustomTaskCount    int64               `protobuf:"zigzag64,5,opt,name=customTaskCount,proto3" json:"customTaskCount,omitempty" yaml:"customTaskCount"`
	AuxiliaryTaskCount int64               `protobuf:"zigzag64,6,opt,name=auxiliaryTaskCount,proto3" json:"auxiliaryTaskCount,omitempty" yaml:"auxiliaryTaskCount"`
}

type ProjectCluster struct {
	ProjectID string `protobuf:"bytes,1,opt,name=projectID,proto3" json:"projectID,omitempty" yaml:"projectID"`
	ClusterID string `protobuf:"bytes,2,opt,name=clusterID,proto3" json:"clusterID,omitempty" yaml:"clusterID"`
}

type ClusterCondition struct {
	Type        ClusterConditionType `protobuf:"varint,1,opt,name=type,proto3,enum=rafay.dev.types.infra.v3.ClusterConditionType" json:"type,omitempty" yaml:"type"`
	Status      RafayConditionStatus `protobuf:"varint,2,opt,name=status,proto3,enum=rafay.dev.types.common.v3.RafayConditionStatus" json:"status,omitempty" yaml:"status"`
	LastUpdated time.Time            `protobuf:"bytes,3,opt,name=lastUpdated,proto3" json:"lastUpdated,omitempty" yaml:"lastUpdated"`
	Reason      string               `protobuf:"bytes,4,opt,name=reason,proto3" json:"reason,omitempty" yaml:"reason"`
}

// model for bootstrap file download
type BootstrapFileDownload struct {
	ContentType string `json:"content_type" yaml:"content_type"`
	Data        string `json:"data" yaml:"data"`
}

type ClusterDisplayDetails struct {
	Name        string `json:"name,omitempty" yaml:"name"`
	ClusterType string `json:"cluster_type,omitempty" yaml:"cluster_type"`
}

type ClusterList struct {
	ApiVersion string       `protobuf:"bytes,1,opt,name=apiVersion,proto3" json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string       `protobuf:"bytes,2,opt,name=kind,proto3" json:"kind,omitempty" yaml:"kind"`
	Metadata   ListMetadata `protobuf:"bytes,3,opt,name=metadata,proto3" json:"metadata,omitempty" yaml:"metadata"`
	Items      []*Cluster   `protobuf:"bytes,4,rep,name=items,proto3" json:"items,omitempty" yaml:"items"`
}

type ProvisionParams struct {
	EnvironmentProvider  string `protobuf:"bytes,1,opt,name=environmentProvider,proto3" json:"environmentProvider,omitempty" yaml:"environmentProvider"`
	KubernetesProvider   string `protobuf:"bytes,2,opt,name=kubernetesProvider,proto3" json:"kubernetesProvider,omitempty" yaml:"kubernetesProvider"`
	ProvisionEnvironment string `protobuf:"bytes,3,opt,name=provisionEnvironment,proto3" json:"provisionEnvironment,omitempty" yaml:"provisionEnvironment"`
	ProvisionPackageType string `protobuf:"bytes,4,opt,name=provisionPackageType,proto3" json:"provisionPackageType,omitempty" yaml:"provisionPackageType"`
	ProvisionType        string `protobuf:"bytes,5,opt,name=provisionType,proto3" json:"provisionType,omitempty" yaml:"provisionType"`
	State                string `protobuf:"bytes,6,opt,name=state,proto3" json:"state,omitempty" yaml:"state"`
}

type ClusterShareMode int32

const (
	ClusterShareMode_ClusterShareModeNotSet ClusterShareMode = 0
	ClusterShareMode_ALL                    ClusterShareMode = 1
	ClusterShareMode_CUSTOM                 ClusterShareMode = 2
)

// Enum value maps for ClusterShareMode.
var (
	ClusterShareMode_name = map[int32]string{
		0: "ClusterShareModeNotSet",
		1: "ALL",
		2: "CUSTOM",
	}
	ClusterShareMode_value = map[string]int32{
		"ClusterShareModeNotSet": 0,
		"ALL":                    1,
		"CUSTOM":                 2,
	}
)

func (x ClusterShareMode) Enum() *ClusterShareMode {
	p := new(ClusterShareMode)
	*p = x
	return p
}

type ProxyConfig struct {
	HttpProxy              string `protobuf:"bytes,1,opt,name=httpProxy,proto3" json:"httpProxy,omitempty"`
	HttpsProxy             string `protobuf:"bytes,2,opt,name=httpsProxy,proto3" json:"httpsProxy,omitempty"`
	NoProxy                string `protobuf:"bytes,3,opt,name=noProxy,proto3" json:"noProxy,omitempty"`
	ProxyAuth              string `protobuf:"bytes,4,opt,name=proxyAuth,proto3" json:"proxyAuth,omitempty"`
	AllowInsecureBootstrap bool   `protobuf:"varint,5,opt,name=allowInsecureBootstrap,proto3" json:"allowInsecureBootstrap,omitempty"`
	Enabled                bool   `protobuf:"varint,6,opt,name=enabled,proto3" json:"enabled,omitempty"`
	BootstrapCA            string `protobuf:"bytes,7,opt,name=bootstrapCA,proto3" json:"bootstrapCA,omitempty"`
}

type Health int32

const (
	Health_EDGE_IGNORE       Health = 0
	Health_EDGE_HEALTHY      Health = 1
	Health_EDGE_UNHEALTHY    Health = 2
	Health_EDGE_DISCONNECTED Health = 3
)

// Enum value maps for Health.
var (
	Health_name = map[int32]string{
		0: "EDGE_IGNORE",
		1: "EDGE_HEALTHY",
		2: "EDGE_UNHEALTHY",
		3: "EDGE_DISCONNECTED",
	}
	Health_value = map[string]int32{
		"EDGE_IGNORE":       0,
		"EDGE_HEALTHY":      1,
		"EDGE_UNHEALTHY":    2,
		"EDGE_DISCONNECTED": 3,
	}
)

func (x Health) Enum() *Health {
	p := new(Health)
	*p = x
	return p
}

type ClusterConditionType int32

const (
	ClusterConditionType_ClusterBlueprintSync     ClusterConditionType = 0
	ClusterConditionType_ClusterApprove           ClusterConditionType = 1
	ClusterConditionType_ClusterCheckIn           ClusterConditionType = 2
	ClusterConditionType_ClusterNodeSync          ClusterConditionType = 3
	ClusterConditionType_ClusterRegister          ClusterConditionType = 4
	ClusterConditionType_ClusterNamespaceSync     ClusterConditionType = 5
	ClusterConditionType_ClusterReady             ClusterConditionType = 6
	ClusterConditionType_ClusterAuxiliaryTaskSync ClusterConditionType = 7
	ClusterConditionType_ClusterBootstrapAgent    ClusterConditionType = 8
	ClusterConditionType_ClusterDelete            ClusterConditionType = 9
)

// Enum value maps for ClusterConditionType.
var (
	ClusterConditionType_name = map[int32]string{
		0: "ClusterBlueprintSync",
		1: "ClusterApprove",
		2: "ClusterCheckIn",
		3: "ClusterNodeSync",
		4: "ClusterRegister",
		5: "ClusterNamespaceSync",
		6: "ClusterReady",
		7: "ClusterAuxiliaryTaskSync",
		8: "ClusterBootstrapAgent",
		9: "ClusterDelete",
	}
	ClusterConditionType_value = map[string]int32{
		"ClusterBlueprintSync":     0,
		"ClusterApprove":           1,
		"ClusterCheckIn":           2,
		"ClusterNodeSync":          3,
		"ClusterRegister":          4,
		"ClusterNamespaceSync":     5,
		"ClusterReady":             6,
		"ClusterAuxiliaryTaskSync": 7,
		"ClusterBootstrapAgent":    8,
		"ClusterDelete":            9,
	}
)

func (x ClusterConditionType) Enum() *ClusterConditionType {
	p := new(ClusterConditionType)
	*p = x
	return p
}

// RafayConditionStatus is the status of the status condition
type RafayConditionStatus int32

const (
	RafayConditionStatus_NotSet     RafayConditionStatus = 0
	RafayConditionStatus_Pending    RafayConditionStatus = 1
	RafayConditionStatus_InProgress RafayConditionStatus = 2
	RafayConditionStatus_Success    RafayConditionStatus = 3
	RafayConditionStatus_Failed     RafayConditionStatus = 4
	RafayConditionStatus_Retry      RafayConditionStatus = 5
	RafayConditionStatus_Skipped    RafayConditionStatus = 6
	RafayConditionStatus_Stopped    RafayConditionStatus = 7
	RafayConditionStatus_Expired    RafayConditionStatus = 8
	RafayConditionStatus_Stopping   RafayConditionStatus = 9
	RafayConditionStatus_Submitted  RafayConditionStatus = 10
)

// Enum value maps for RafayConditionStatus.
var (
	RafayConditionStatus_name = map[int32]string{
		0:  "NotSet",
		1:  "Pending",
		2:  "InProgress",
		3:  "Success",
		4:  "Failed",
		5:  "Retry",
		6:  "Skipped",
		7:  "Stopped",
		8:  "Expired",
		9:  "Stopping",
		10: "Submitted",
	}
	RafayConditionStatus_value = map[string]int32{
		"NotSet":     0,
		"Pending":    1,
		"InProgress": 2,
		"Success":    3,
		"Failed":     4,
		"Retry":      5,
		"Skipped":    6,
		"Stopped":    7,
		"Expired":    8,
		"Stopping":   9,
		"Submitted":  10,
	}
)

func (x RafayConditionStatus) Enum() *RafayConditionStatus {
	p := new(RafayConditionStatus)
	*p = x
	return p
}
