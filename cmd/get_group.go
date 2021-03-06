package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetGroupCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "group",
		Aliases: []string{"gr", "group"},
		Short:   "Get a list of groups or a single group",
		Long:    `Retrieves a list of groups or a single group.`,
		Example: `
   List all groups
     pctl get groups

   Show more about a specific groups
     pctl get group <group-name>
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)
	return cmd
}
