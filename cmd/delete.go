package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/log"
	"github.com/spf13/cobra"
)

func newDeleteCmd(logger log.Logger, config *config.Config) *cobra.Command {
	// cmd represents the delete command
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete various resources in Console",
		Long:  `Delete clusters, users, groups and other resources in your current project`,
	}

	// add subcommands here
	cmd.AddCommand(
		newDeleteClusterCmd(commands.NewDeleteClusterOptions(logger)),
		newDeleteProjectCmd(commands.NewDeleteProjectOptions(logger, config)),
		newDeleteGroupCmd(commands.NewDeleteGroupOptions(logger, config)),
		newDeleteUserCmd(commands.NewDeleteUserOptions(logger)),
		newDeleteOIDCProviderCmd(commands.NewDeleteOIDCProviderOptions(logger, config)),
		newDeleteRoleCmd(commands.NewDeleteRoleOptions(logger, config)),
	)

	return cmd
}
