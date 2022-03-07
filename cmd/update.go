package cmd

import (
	"github.com/spf13/cobra"

	"github.com/RafaySystems/rcloud-cli/pkg/commands"
	"github.com/RafaySystems/rcloud-cli/pkg/config"
	"github.com/RafaySystems/rcloud-cli/pkg/log"
)

func newUpdateCmd(logger log.Logger, config *config.Config) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update the details of various resource",
		Long:  `update the details of various resource such as clusters, users, groups etc.`,
	}

	// add subcommands
	cmd.AddCommand(
		newUpdateGroupassociationCmd(commands.NewUpdateGroupassociationOptions(logger, config)),
	)
	return cmd
}
