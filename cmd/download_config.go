package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newConfigDownloadCmd(o commands.CmdOptions) *cobra.Command {

	// cmd represents the get command
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download CLI config.",
		Long: `Download CLI config.
Examples:
	# Download a CLI config from http://console.paralus.local (provide user credentials when asked)
	pctl config download http://console.paralus.local

	# Download a CLI config using provided user credentials.
	pctl config download https://console.paralus-host.com --email myemail@host.com --password $PARALUS_PASSWORD

	# Download a CLI config to $HOME/path/to/cli.config
	pctl config download http://console.paralus.local --to-file $HOME/path/to/cli.config`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
