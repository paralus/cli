package output

import (
	"github.com/spf13/cobra"

	"github.com/RafaySystems/rcloud-cli/pkg/log"
)

func NewLocationListSpec(obj interface{}) *OutputListSpec {
	countInt := GetCount(obj)
	log.GetLogger().Debugf("count = %d", countInt)

	columnSpec := []ColumnSpec{
		{Header: "NAME", JsonPath: "name"},
		{Header: "CITY", JsonPath: "city"},
		{Header: "STATE", JsonPath: "state"},
		{Header: "COUNTRY", JsonPath: "country"},
		{Header: "LATITUDE", JsonPath: "latitude"},
		{Header: "LONGITUDE", JsonPath: "longitude"},
	}

	spec := &OutputListSpec{
		NRows:   countInt,
		Base:    "$.results[%d]",
		Columns: columnSpec,
	}

	return spec
}

func LocationListPrint(cmd *cobra.Command, jsonObj []byte) error {
	return PrintList(cmd, jsonObj, NewLocationListSpec)

}
