package cmd

import (
	"github.com/spf13/cobra"

	"github.com/paralus/cli/pkg/commands"
	"github.com/paralus/cli/pkg/log"
)

func newCreateClusterCmd(importedOption commands.CmdOptions, logger log.Logger) *cobra.Command {
	// cmd represents the cluster command
	cmd := &cobra.Command{
		Use:     "cluster",
		Aliases: []string{"c"},
		Short:   "Create or import a cluster",
		Long:    `Provision a cluster or import an existing cluster.`,
		Example: `
  Using config file:
    pctl create cluster -f cluster-config.yml 
`,
		Args: importedOption.Validate,
		RunE: importedOption.Run,
	}

	importedOption.AddFlags(cmd)

	cmd.AddCommand(newCreateClusterImportCmd(importedOption))

	return cmd

}
