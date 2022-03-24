package cmd

import (
	"github.com/rafaylabs/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newDeleteUserCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the create command
	cmd := &cobra.Command{
		Use:     "user <user-name>",
		Aliases: []string{"u"},
		Short:   "Delete a new user",
		Long:    "Delete a new user",
		Example: `
Using command:
	rctl delete user john.doe@example.com 

`,

		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
