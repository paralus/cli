package cmd

import (
	"github.com/RafaySystems/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetProjectCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"p", "projects"},
		Short:   "Get a list of projects or a single project",
		Long:    `Retrieves a list of projects or a single project.`,
		Example: `
   List all projects
     rctl get projects

   Show more about a specific project
     rctl get project <project-name>
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)
	return cmd
}
