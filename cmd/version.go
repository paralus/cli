package cmd

import (
	"github.com/paralus/cli/pkg/versioninfo"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays version of the CLI utility",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("%v", versioninfo.Get())
		},
	}

	return cmd
}
