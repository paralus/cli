package cmd

import (
	"github.com/RafayLabs/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newUpdateClusterCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "cluster <cluster-name>",
		Aliases: []string{"c", "clusters"},
		Short:   "Update a cluster",
		Long:    `Update a cluster's description, labels, etc.`,
		Args:    o.Validate,
		RunE:    o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
