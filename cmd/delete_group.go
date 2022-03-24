package cmd

import (
	"github.com/rafaylabs/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newDeleteGroupCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:     "group <group>",
		Aliases: []string{"r"},
		Short:   "Delete a group",
		Long:    "Delete a group by name, or multiple groups by entering the names in a space-separated list",
		Args:    o.Validate,
		RunE:    o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
