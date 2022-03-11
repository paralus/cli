package cmd

import (
	"github.com/RafaySystems/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newDownloadKubeconfigCmd(o commands.CmdOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kubeconfig",
		Short:   "Download the generated kubeconfig",
		Long:    "Download the generated kubeconfig",
		Aliases: []string{"k", "kc"},
		Args:    o.Validate,
		RunE:    o.Run,
	}

	o.AddFlags(cmd)

	return cmd
}
