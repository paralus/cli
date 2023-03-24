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
		pctl config download http://console.paralus.local (value of host endpoint of paralus dashboard)`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
