package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newCreateProjectCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the create command
	cmd := &cobra.Command{
		Use:     "project <project-name>",
		Aliases: []string{"p"},
		Short:   "Create a new project",
		Long:    "Create a new project",
		Example: `
Using command:
	pctl create p project1 --desc "Description of the project"
Using file:
	pctl create p -f <path-to-namespace-yaml> --v3
	`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
