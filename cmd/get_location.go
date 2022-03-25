package cmd

import (
	"github.com/RafayLabs/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newGetLocationCmd(o commands.CmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "location",
		Aliases: []string{"l", "locations"},
		Short:   "Get locations or a location",
		Long:    `Get locations or a location`,
		Args:    o.Validate,
		RunE:    o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
