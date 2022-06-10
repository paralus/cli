package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetClusterBootstrapCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "clusterbootstrap",
		Aliases: []string{"cb"},
		Short:   "Get a cluster bootstrap for clusters of type import",
		Long:    `Get a cluster bootstrap for clusters of type import.`,
		Example: "pctl get cb <cluster-name> | kubectl apply -f -",
		Args:    o.Validate,
		RunE:    o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
