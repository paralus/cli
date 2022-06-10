package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newDeleteClusterCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:     "cluster",
		Aliases: []string{"c"},
		Short:   "Delete a cluster",
		Long:    `Delete a cluster`,
		Example: `
  Using command(s):
	pctl delete cluster <cluster1> [<cluster2> ...]
	pctl delete cluster <cluster1> [<cluster2> ...] -t gke --gkeproject <project-name> --zone <compute-zone> --region <compute-region>
  `,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
