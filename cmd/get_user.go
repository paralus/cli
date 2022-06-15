package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetUserCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "user",
		Aliases: []string{"u", "users"},
		Short:   "Get a list of users or a single user",
		Long:    `Retrieves a list of users or a single user.`,
		Example: `
   List all user
     pctl get users

   Show more about a specific groups
     pctl get user <user_name>
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)
	return cmd
}
