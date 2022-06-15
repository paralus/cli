package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newCreateUserCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the create command
	cmd := &cobra.Command{
		Use:     "user <group-name>",
		Aliases: []string{"u"},
		Short:   "Create a new user",
		Long:    "Create a new user",
		Example: `
Using command:
	pctl create user john.doe@example.com 
	pctl create user john.doe@example.com --console John, Doe
	pctl create user john.doe@example.com  --groups testingGroup, productionGroup --console John, Doe, 4089382091
`,

		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
