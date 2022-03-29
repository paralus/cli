package cmd

import (
	"os"

	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newCompletionCmd(*zap.SugaredLogger, *config.Config) *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(rctl completion bash)

# To load completions for each session, execute once:
Linux:
  $ rctl completion bash > /etc/bash_completion.d/rctl
MacOS:
  $ rctl completion bash > /usr/local/etc/bash_completion.d/rctl

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ rctl completion zsh > "${fpath[1]}/_rctl"

# You will need to start a new shell for this setup to take effect.

Fish:

$ rctl completion fish | source

# To load completions for each session, execute once:
$ rctl completion fish > ~/.config/fish/completions/rctl.fish
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletion(os.Stdout)
			}
		},
	}
	return completionCmd
}
