package cmd

import (
	"github.com/RafaySystems/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newDeleteIdpCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:   "oidc <oidc-provider-name>",
		Short: "Delete a OIDC",
		Long:  "Delete a OIDC by name, or multiple OIDC's by entering the names in a space-separated list",
		Args:  o.Validate,
		RunE:  o.Run,
		Example: `
Delete Identify Provider
	rctl delete oidc name1

Delete Identify Provider(s)
	rctl delete oidc name1 name2 name13 
`,
	}

	o.AddFlags(cmd)

	return cmd
}
