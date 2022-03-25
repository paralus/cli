package cmd

import (
	"github.com/RafayLabs/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newCreateRoleCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the create command
	cmd := &cobra.Command{
		Use:   "role <custom-role-name>",
		Short: "Create a custom role",
		Long:  "Create a custom role with varying set of permissions",
		Example: `

Basic Info : This command is used to create custom roles
----------------------
custom-role-name : Provide a custom role name

permissions: List of all the role permissions associated to this role
----------------------

Examples:
	Basic Command :
		rctl create role clusterview --permissions project.read,cluster.read,project.clusterctl.read
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
