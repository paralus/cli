package cmd

import (
	"github.com/RafayLabs/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetOIDCProviderCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:   "oidc",
		Short: "Get list of all OIDC names, domains and group attribute names",
		Long:  "Get list of all OIDC names, domains and group attribute names",
		Example: `
 List specific Identify Provider details
     rctl get oidc <oidcname>

List all Identify Providers details
     rctl get oidc
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)
	return cmd
}
