package cluster

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/constants"
	"github.com/paralus/cli/pkg/rerror"
	commonv3 "github.com/paralus/paralus/proto/types/commonpb/v3"
	infrav3 "github.com/paralus/paralus/proto/types/infrapb/v3"
)

// NewImportCluster will create a new cluster of type import
func NewImportCluster(name, location, project string) (string, error) {
	importCluster := infrav3.Cluster{
		Kind: "Cluster",
		Metadata: &commonv3.Metadata{
			Name:    name,
			Project: project,
		},
		Spec: &infrav3.ClusterSpec{
			Metro: &infrav3.Metro{
				Name: location,
			},
			ClusterType: constants.CLUSTER_TYPE_IMPORT,
		},
	}

	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster", project)
	resp, err := auth.AuthAndRequest(uri, "POST", importCluster)

	if err != nil {
		return resp, err
	}
	return resp, nil
}

// NewImportClusterOpenshift will create a new cluster of type import
func NewImportClusterOpenshift(name, location, project string) (string, error) {
	importCluster := infrav3.Cluster{
		Kind: "Cluster",
		Metadata: &commonv3.Metadata{
			Name:    name,
			Project: project,
		},
		Spec: &infrav3.ClusterSpec{
			Metro: &infrav3.Metro{
				Name: location,
			},
			ClusterType: constants.CLUSTER_TYPE_IMPORT,
			Params: &infrav3.ProvisionParams{
				EnvironmentProvider:  "",
				KubernetesProvider:   "OPENSHIFT",
				ProvisionEnvironment: "ONPREM",
				ProvisionPackageType: "",
				ProvisionType:        "IMPORT",
				State:                "PROVISION",
			},
		},
	}

	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster", project)
	resp, err := auth.AuthAndRequest(uri, "POST", importCluster)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// NewImportClusterAKS will create a new cluster of type import
func NewImportClusterAKS(name, location, project string) (string, error) {

	importCluster := infrav3.Cluster{
		Kind: "Cluster",
		Metadata: &commonv3.Metadata{
			Name:    name,
			Project: project,
		},
		Spec: &infrav3.ClusterSpec{
			Metro: &infrav3.Metro{
				Name: location,
			},
			ClusterType: constants.CLUSTER_TYPE_IMPORT,
			Params: &infrav3.ProvisionParams{
				EnvironmentProvider:  "AZURE",
				KubernetesProvider:   "AKS",
				ProvisionEnvironment: "CLOUD",
				ProvisionPackageType: "",
				ProvisionType:        "IMPORT",
				State:                "PROVISION",
			},
		},
	}

	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster", project)
	resp, err := auth.AuthAndRequest(uri, "POST", importCluster)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// NewImportClusterGKE will create a new cluster of type import
func NewImportClusterGKE(name, location, project string) (string, error) {

	importCluster := infrav3.Cluster{
		Kind: "Cluster",
		Metadata: &commonv3.Metadata{
			Name:    name,
			Project: project,
		},
		Spec: &infrav3.ClusterSpec{
			Metro: &infrav3.Metro{
				Name: location,
			},
			ClusterType: constants.CLUSTER_TYPE_IMPORT,
			Params: &infrav3.ProvisionParams{
				EnvironmentProvider:  "GCP",
				KubernetesProvider:   "GKE",
				ProvisionEnvironment: "CLOUD",
				ProvisionPackageType: "",
				ProvisionType:        "IMPORT",
				State:                "PROVISION",
			},
		},
	}

	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster", project)
	resp, err := auth.AuthAndRequest(uri, "POST", importCluster)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func NewImportClusterMKS(name, location, project, K8Sversion, OsVersion, defaultStorageClass string, Storageclassmap map[string]string) (string, error) {

	importCluster := infrav3.Cluster{
		Kind: "Cluster",
		Metadata: &commonv3.Metadata{
			Name:    name,
			Project: project,
		},
		Spec: &infrav3.ClusterSpec{
			Metro: &infrav3.Metro{
				Name: location,
			},
			ClusterType: constants.CLUSTER_TYPE_IMPORT,
			Params: &infrav3.ProvisionParams{
				EnvironmentProvider:  "",
				KubernetesProvider:   "MKS",
				ProvisionEnvironment: "ONPREM",
				ProvisionPackageType: "LINUX",
				ProvisionType:        "IMPORT",
				State:                "CONFIG",
			},
		},
	}

	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster", project)
	resp, err := auth.AuthAndRequest(uri, "POST", importCluster)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// ListAllClusters uses the lower level func ListClusters to retrieve a list of all clusters
func ListAllClusters(projectId string) ([]*infrav3.Cluster, error) {
	var clusters []*infrav3.Cluster
	limit := 10000
	c, count, err := ListClusters(projectId, limit, 0)
	if err != nil {
		return nil, err
	}
	clusters = c
	for count > limit {
		offset := limit
		limit = count
		c, _, err = ListClusters(projectId, limit, offset)
		if err != nil {
			return clusters, err
		}
		clusters = append(clusters, c...)
	}
	return clusters, nil
}

/*
ListClusters paginates through a list of clusters
*/
func ListClusters(project string, limit, offset int) ([]*infrav3.Cluster, int, error) {
	// check to make sure the limit or offset is not negative
	if limit < 0 || offset < 0 {
		return nil, 0, fmt.Errorf("provided limit (%d) or offset (%d) cannot be negative", limit, offset)
	}
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster?limit=%d&offset=%d", project, limit, offset)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, 0, rerror.CrudErr{
			Type: "cluster",
			Name: "",
			Op:   "list",
		}
	}
	a := infrav3.ClusterList{}
	_ = json.Unmarshal([]byte(resp), &a)
	return a.Items, int(a.Metadata.Count), nil
}

func getClusterFast(name, project string) (*infrav3.Cluster, error) {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster/%s", project, name)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, errors.New("error fetching cluster details")
	}
	var cluster infrav3.Cluster
	if err := json.Unmarshal([]byte(resp), &cluster); err != nil {
		return nil, errors.New("error unmarshalling cluster details")
	}
	return &cluster, nil
}

/*
GetCluster gets an cluster based on the name provided. It calls ListAllClusters, and scan through the list
for the name of the cluster. Returns nil if such addon does not exist, or returns an error if there was
and error fetching all of the addons
*/
func GetCluster(name, project string) (*infrav3.Cluster, error) {

	// first try using the name filter
	cluster, err := getClusterFast(name, project)
	if err == nil {
		return cluster, nil
	}

	// get list of clusters
	c, err := ListAllClusters(project)
	if err != nil {
		return nil, err
	}
	for _, a := range c {
		if a.Metadata.Name == name {
			return a, nil
		}
	}

	return nil, rerror.ResourceNotFound{
		Type: "cluster",
		Name: name,
	}
}

func DeleteCluster(name, project string) error {
	// get cluster
	_, err := GetCluster(name, project)
	if err != nil {
		return err
	}

	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster/%s", project, name)
	_, err = auth.AuthAndRequest(uri, "DELETE", nil)
	if err != nil {
		return rerror.CrudErr{
			Type: "cluster",
			Name: "",
			Op:   "delete",
		}
	}

	return nil
}

// GetBootstrapFile will retrieve the bootstrap file for imported clusters
func GetBootstrapFile(name, project string) (string, error) {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster/%s/download", project, name)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return "", rerror.CrudErr{
			Type: "cluster bootstrap",
			Name: name,
			Op:   "get",
		}
	}

	f := &commonv3.HttpBody{}
	err = json.Unmarshal([]byte(resp), f)
	if err != nil {
		return "", err
	}

	return string(f.Data), nil
}

// Update cluster takes the updated cluster details and sends it to the core
func UpdateCluster(cluster *infrav3.Cluster) error {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster/%s", cluster.Metadata.Project, cluster.Metadata.Name)
	_, err := auth.AuthAndRequest(uri, "PUT", cluster)
	if err != nil {
		return rerror.CrudErr{
			Type: "cluster",
			Name: cluster.Metadata.Name,
			Op:   "update",
		}
	}
	return nil
}

// Update cluster takes the updated cluster details and sends it to the core
func CreateCluster(cluster *infrav3.Cluster) error {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster", cluster.Metadata.Project)
	_, err := auth.AuthAndRequest(uri, "POST", cluster)
	if err != nil {
		return rerror.CrudErr{
			Type: "cluster",
			Name: cluster.Metadata.Name,
			Op:   "create",
		}
	}
	return nil
}
