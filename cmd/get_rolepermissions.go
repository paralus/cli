package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetRolePermissionCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "rolepermissions",
		Aliases: []string{"rp", "rolepermission"},
		Short:   "Get a list of role permissions",
		Long:    `Retrieves a list of roles permissions.`,
		Args:    o.Validate,
		RunE:    o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
