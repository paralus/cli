package cmd

import (
	"github.com/spf13/cobra"

	"github.com/RafayLabs/rcloud-cli/pkg/commands"
)

func newApplyCmd(applyOption commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply the config of various resource",
		Long:  `Apply configuration changes on various resources like clusters, users, groups etc.`,
		Example: `
  Using config file:
    rctl apply -f cluster-config.yml 
`,
		Args: applyOption.Validate,
		RunE: applyOption.Run,
	}

	return cmd
}
