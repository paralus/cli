package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newDeleteIdpCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:   "idp <idp-name>",
		Short: "Delete a Idp",
		Long:  "Delete a Idp by name, or multiple IDp's by entering the names in a space-separated list",
		Args:  o.Validate,
		RunE:  o.Run,
		Example: `
Delete Identify Provider
	pctl delete idp name1

Delete Identify Provider(s)
	pctl delete idp name1 name2 name13 

`,
	}

	o.AddFlags(cmd)

	return cmd
}
