package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetRoleCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "roles",
		Aliases: []string{"r", "role"},
		Short:   "Get a list of roles",
		Long:    `Retrieves a list of roles.`,
		Example: `
   List all roles
     pctl get roles

   Show more about a specific role
     pctl get role <role-name>
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
