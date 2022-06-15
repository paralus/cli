package commands

import (
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/test"
	"github.com/spf13/cobra"
)

func newTestWithConfigCmd(Usage string, ArgsFunc, RunEFunc func(cmd *cobra.Command, args []string) error) *cobra.Command {
	configRunFunc := func(cmd *cobra.Command, args []string) error {
		err := setUpConfigBeforeTest()
		if err != nil {
			return err
		}
		return RunEFunc(cmd, args)
	}
	return newTestCmd(Usage, ArgsFunc, configRunFunc)
}

func newTestCmd(Usage string, ArgsFunc, RunEFunc func(cmd *cobra.Command, args []string) error) *cobra.Command {
	r := cobra.Command{Use: Usage, Args: ArgsFunc, RunE: func(cmd *cobra.Command, args []string) error {
		return RunEFunc(cmd, args)
	}}
	NewGlobalOptions(test.NewNoopLogger(), config.GetConfig()).AddFlags(&r)
	return &r
}

func setUpConfigBeforeTest() error {
	return test.SetUpConfigBeforeTestWithConfigPath("testdata", "config.json")
}
