package cmd

import (
	"github.com/RafaySystems/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newApplyResourceCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the cluster command
	cmd := &cobra.Command{
		Use:   "-f <config.yml>",
		Short: "Apply configuration changes",
		Long:  `Apply configuration changes on various resources like clusters, users, groups etc.`,
		Example: `
  Using config file:
    rctl apply -f cluster-config.yml 
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
