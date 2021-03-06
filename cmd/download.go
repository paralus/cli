package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/log"
	"github.com/spf13/cobra"
)

func newDownloadCmd(logger log.Logger, config *config.Config) *cobra.Command {
	// createCmd represents the create command
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download various resources in Console",
		Long:  `Download various resources in Console`,
	}

	// add subcommands here
	cmd.AddCommand(
		newDownloadKubeconfigCmd(commands.NewDownloadKubeconfigOptions(logger)),
	)

	return cmd
}
