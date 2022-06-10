package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetIdpCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:   "idp",
		Short: "Get list of all Idp names, domains and group attribute names",
		Long:  "Get list of all Idp names, domains and group attribute names",
		Example: `
 List specific Identify Provider details
     pctl get idp <idpname>

List all Identify Providers detailss
     pctl get idp
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)
	return cmd
}
