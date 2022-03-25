package output

import (
	"github.com/RafayLabs/rcloud-cli/pkg/log"

	"github.com/oliveagle/jsonpath"
	"github.com/spf13/cobra"
)

func GetCount(obj interface{}) int {
	countInt := 0
	count, err := jsonpath.JsonPathLookup(obj, "$.count")
	if err == nil {
		if countF, ok := count.(float64); ok != true {
			log.GetLogger().Infof("Failed to convert 'count' into float64")
			countInt = 0
		} else {
			countInt = int(countF)
		}
	} else {
		log.GetLogger().Debugf("Failed to get 'count' from result %s", err)
	}

	return countInt
}

func NewClusterListSpec(obj interface{}) *OutputListSpec {
	countInt := GetCount(obj)

	log.GetLogger().Debugf("count = %d", countInt)
	columnSpec := []ColumnSpec{
		{Header: "NAME", JsonPath: "name"},
		{Header: "TYPE", JsonPath: "cluster_type"},
		{Header: "LOCATION", JsonPath: "Metro.name"},
		{Header: "STATUS", JsonPath: "status"},
	}

	spec := &OutputListSpec{
		NRows:   countInt,
		Base:    "$.results[%d]",
		Columns: columnSpec,
	}

	return spec
}

func NewClusterObjectSpec(obj interface{}, clusterName string) *OutputObjectSpec {
	countInt := GetCount(obj)

	log.GetLogger().Debugf("count = %d", countInt)
	matchSpec := MatchSpec{JsonPath: "name", Value: clusterName}

	spec := &OutputObjectSpec{
		NRows: countInt,
		Base:  "$.results[%d]",
		Match: matchSpec,
	}

	return spec
}

func ClusterListPrint(cmd *cobra.Command, jsonObj []byte) error {
	return PrintList(cmd, jsonObj, NewClusterListSpec)
}

func ClusterSelectPrint(cmd *cobra.Command, jsonObj []byte, clusterName string) error {
	return PrintSelectObject(cmd, jsonObj, clusterName, NewClusterObjectSpec)
}
