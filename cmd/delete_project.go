package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newDeleteProjectCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:     "project <project>",
		Aliases: []string{"p"},
		Short:   "Delete a project",
		Long:    "Delete a project by name, or multiple projects by entering the names in a space-separated list",
		Example: `
Using command:
	pctl delete p project1
Using file:
	pctl delete p -f <path-to-namespace-yaml> --v3
	`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
