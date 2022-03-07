package cmd

import (
	"fmt"
	"os"

	"github.com/RafaySystems/rcloud-cli/pkg/commands"
	"github.com/RafaySystems/rcloud-cli/pkg/config"
	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"github.com/RafaySystems/rcloud-cli/pkg/output"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func newRootCmd() *cobra.Command {
	c := config.GetConfig()
	logger := log.GetLogger()
	o := commands.NewGlobalOptions(logger, c)
	// this cmd represents the base command when called without any subcommands
	cmd := &cobra.Command{
		Use:               "rctl",
		Short:             "A CLI tool to manage resources.",
		Long:              `A CLI tool to manage resources.`,
		TraverseChildren:  true,
		SilenceUsage:      true,
		PersistentPreRunE: o.Run,
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			/* Display and exit. */
			output.Exit()
		},
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	o.AddFlags(cmd)

	// add subcommands here
	cmd.AddCommand(
		// new commands
		newCreateCmd(logger, c),
		newGetCmd(logger, c),
		newUpdateCmd(logger, c),
		newDeleteCmd(logger, c),
		newApplyCmd(commands.NewApplyResourcesOptions(logger, c)),

		KubeconfigCmd,
	)

	cmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		if err := c.Help(); err != nil {
			return err
		}
		return err
	})

	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// if any of the subcommands run into an error, it shows up here
func Execute() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// run when each command's execute method is called
	// do command wide inits here
	cobra.OnInitialize(initConfig)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".rctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.GetLogger().Infof("Using config file:", viper.ConfigFileUsed())
	}
}
