package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newCreateGroupCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the create command
	cmd := &cobra.Command{
		Use:     "group <group-name>",
		Aliases: []string{"r"},
		Short:   "Create a new group",
		Long:    "Create a new group",
		Example: `
Using command:
	pctl create group sample-group --desc "Description of the group"`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
