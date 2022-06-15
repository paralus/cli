package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newDeleteRoleCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:   "role <custom-role-name>",
		Short: "Delete a Role",
		Long:  "Delete a Custom Role",
		Args:  o.Validate,
		RunE:  o.Run,
		Example: `
Delete Role
	pctl delete role customrole

Delete Role(s)
	pctl delete role name1 name2 name3
`,
	}

	o.AddFlags(cmd)

	return cmd
}
