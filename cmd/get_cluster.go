package cmd

import (
	"github.com/RafayLabs/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetClusterCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "cluster",
		Aliases: []string{"c", "clusters"},
		Short:   "Get a list of clusters or a single cluster",
		Long:    `Retrieves a list of clusters or a single cluster.`,
		Args:    o.Validate,
		RunE:    o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
