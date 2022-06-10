package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/log"
	"github.com/spf13/cobra"
)

func newGetCmd(logger log.Logger, config *config.Config) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get various resources in Console",
		Long:  `get clusters, users, groups and other resources in your current project`,
	}

	// add subcommands
	cmd.AddCommand(
		newGetClusterCmd(
			commands.NewGetClusterOptions(logger),
		),
		newGetClusterBootstrapCmd(commands.NewGetClusterBootstrapOptions(logger)),
		newGetProjectCmd(commands.NewGetProjectOptions(logger)),
		newGetLocationCmd(commands.NewGetLocationOptions(logger)),
		newGetUserCmd(commands.NewGetUserOptions(logger)),
		newGetGroupCmd(commands.NewGetGroupOptions(logger)),
		newGetRoleCmd(commands.NewGetRolesOptions(logger)),
		newGetRolePermissionCmd(commands.NewGetRolePermissionOptions(logger)),
		newGetIdpCmd(commands.NewGetIdpOptions(logger)),
		newGetOIDCProviderCmd(commands.NewGetOIDCProviderOptions(logger)),
	)

	return cmd
}
