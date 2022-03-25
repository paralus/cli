package cmd

import (
	"github.com/spf13/cobra"

	"github.com/RafayLabs/rcloud-cli/pkg/commands"
	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
)

func newCreateCmd(logger log.Logger, config *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create various resources in Console",
		Long:  `Create clusters, namespaces, blueprints, and addons in your current project`,
	}

	// add subcommands here
	cmd.AddCommand(
		newCreateClusterCmd(
			commands.NewCreateClusterImportedOptions(logger),
			logger,
		),
		newCreateProjectCmd(commands.NewCreateProjectOptions(logger, config)),
		newCreateUserCmd(commands.NewCreateUserOptions(logger, config)),
		newCreateGroupCmd(commands.NewCreateGroupOptions(logger, config)),
		newCreateIdpCmd(commands.NewCreateIdpOptions(logger, config)),
		newCreateOIDCProviderCmd(commands.NewCreateOIDCProviderOptions(logger, config)),
		newCreateRoleCmd(commands.NewCreateRoleOptions(logger, config)),
	)

	return cmd
}
