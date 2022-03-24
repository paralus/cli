package commands

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/rafaylabs/rcloud-cli/pkg/cluster"
	"github.com/rafaylabs/rcloud-cli/pkg/config"
	"github.com/rafaylabs/rcloud-cli/pkg/log"
	"github.com/rafaylabs/rcloud-cli/pkg/printer"
	"github.com/spf13/cobra"
)

type GetClusterOptions struct {
	logger log.Logger
}

func NewGetClusterOptions(logger log.Logger) *GetClusterOptions {
	o := new(GetClusterOptions)
	o.logger = logger
	return o
}

func (o *GetClusterOptions) Validate(_ *cobra.Command, _ []string) error {
	return nil
}

func (o *GetClusterOptions) Run(cmd *cobra.Command, args []string) error {
	// get flags
	flags := cmd.Flags()

	// retrieve the project id
	p, err := config.GetProjectIdFromFlagAndConfig(cmd)
	if err != nil {
		return err
	}

	// check if we need to retrieve a singe cluster
	if len(args) != 0 {
		n := args[0]
		cl, err := cluster.GetCluster(n, p)
		if err != nil {
			return err
		}

		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, cl, cmd.OutOrStdout())
			return nil
		}

		// marshall addon to json for table format
		j := make([]string, 1)
		t, _ := json.Marshal(cl)

		j[0] = string(t)
		log.GetLogger().Debug(j[0])

		// set the columns
		c := []string{
			"name",
			"description",
			"project",
			"type",
			"status",
		}
		cp := make(map[string]string)
		cp["name"] = "metadata.name"
		cp["description"] = "metadata.description"
		cp["project"] = "metadata.project"
		cp["type"] = "spec.clusterType"
		cp["status"] = "status"

		printer.PrintTableJsonPath(j, c, cp, os.Stdout)

		return nil
	}

	// list of clusters
	clusters, err := cluster.ListAllClusters(p)
	if err != nil {
		return err
	}

	cm := clusters
	for _, a := range clusters {
		if a.Spec.ClusterType == "imported" {
			if a.Spec.Params != nil {
				if a.Spec.Params.KubernetesProvider == "AKS" {
					a.Spec.ClusterType = "aks"
				} else if a.Spec.Params.KubernetesProvider == "GKE" {
					a.Spec.ClusterType = "gke"
				}
			}
		}
	}

	// json or yaml output
	if flags.Changed("output") {
		o, err := flags.GetString("output")
		if err != nil {
			return err
		}
		printer.PrintOutputJsonOrYaml(o, cm, cmd.OutOrStdout())
		return nil
	}

	// set the columns
	c := []string{
		"name",
		"description",
		"type",
		"ownership",
	}

	rows := make([][]string, 0, len(clusters))
	// set the rows
	for _, c := range clusters {
		rows = append(rows, []string{
			c.Metadata.Name,
			c.Metadata.Description,
			c.Spec.ClusterType,
			c.Metadata.Project,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return strings.Compare(rows[i][0], rows[j][0]) == -1
	})

	printer.PrintTable(c, rows, cmd.OutOrStdout())

	return nil
}

func (o *GetClusterOptions) AddFlags(_ *cobra.Command) {}
