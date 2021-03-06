package cmd

import (
	"github.com/paralus/cli/pkg/commands"
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
     pctl get oidc <oidcname>

List all Identify Providers details
     pctl get oidc
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)
	return cmd
}
