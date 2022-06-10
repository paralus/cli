package cmd

import (
	"os"

	"github.com/paralus/cli/pkg/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newCompletionCmd(*zap.SugaredLogger, *config.Config) *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

$ source <(pctl completion bash)

# To load completions for each session, execute once:
Linux:
  $ pctl completion bash > /etc/bash_completion.d/pctl
MacOS:
  $ pctl completion bash > /usr/local/etc/bash_completion.d/pctl

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ pctl completion zsh > "${fpath[1]}/_pctl"

# You will need to start a new shell for this setup to take effect.

Fish:

$ pctl completion fish | source

# To load completions for each session, execute once:
$ pctl completion fish > ~/.config/fish/completions/pctl.fish
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
