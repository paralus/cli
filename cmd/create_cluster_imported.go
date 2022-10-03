package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newCreateClusterImportCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the cluster command
	cmd := &cobra.Command{
		Use:   "imported <cluster-name>",
		Short: "Import a cluster",
		Long:  `Import an existing cluster.`,
		Example: `
  Using command(s):
    pctl create cluster imported sample-imported-cluster
	
  Using config file:
    pctl create cluster imported -f cluster-config.yml 
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
