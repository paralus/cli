package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/paralus/cli/pkg/commands"
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/context"
	"github.com/paralus/cli/pkg/exit"
	"github.com/paralus/cli/pkg/log"
	"github.com/paralus/cli/pkg/output"
	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration of the CLI utility",
	Long:  "Manage configuration of the CLI utility",
}

var ConfigInitCmd = &cobra.Command{
	Use:   "init",
	Short: "import config",
	Long:  "import config",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		log.GetLogger().Infof("Start [%s %s]", cmd.CommandPath(), filename)
		c := &config.Config{}

		file, err := ioutil.ReadFile(filename)
		if err != nil {
			log.GetLogger().Debugf("Failed to open file %s, [ERROR: %s]", filename, err)
			exit.SetExitWithError(err, fmt.Sprintf("Failed to open file %s", filename))
			return
		}

		err = json.Unmarshal([]byte(file), c)
		if err != nil {
			log.GetLogger().Debugf("Failed to parse file %s, [ERROR: %s]", filename, err)
			exit.SetExitWithError(err, fmt.Sprintf("Failed to parse file as json file %s", filename))
			return
		}

		c.Log(fmt.Sprintf("Init with file %s: ", filename))

		if len(c.Profile) == 0 {
			c.Profile = "prod"
		}

		err = c.MiniCheck()
		if err != nil {
			log.GetLogger().Debugf("Config file has missing fields %s, [ERROR: %s]", filename, err)
			exit.SetExitWithError(err, fmt.Sprintf("Config file has missing fields %s", filename))
			return
		}

		log.GetLogger().Infof("Minicheck ok")

		configFile := context.GetContext().ConfigFilename()
		err = c.Save(configFile)
		if err != nil {
			log.GetLogger().Debugf("Error in saving config file, [ERROR: %s]", err)
			exit.SetExitWithError(err, fmt.Sprintf("Error saving confing file"))
			return
		}

		log.GetLogger().Infof("Save config to file %s", configFile)
		log.GetLogger().Infof("End [%s %s]", cmd.CommandPath(), filename)
	},
}

var ConfigShowCmd = &cobra.Command{
	Use:   "show",
	Short: "display current config",
	Long:  "display current config",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		log.GetLogger().Infof("Start [%s]", cmd.CommandPath())
		showSource, err := cmd.Flags().GetBool("source")
		if err != nil {
			showSource = false
		}
		if showSource {
			output.PrintOutputer(cmd, config.GetConfigTracker())
		} else {
			output.PrintOutputer(cmd, config.GetConfig())
		}
		log.GetLogger().Infof("End [%s]", cmd.CommandPath())
	},
}

var ConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set pctl config",
	Long:  "Set various parameters of pctl config",
}

var ConfigSetProjectCmd = &cobra.Command{
	Use:   "project <project name>",
	Short: "set the project to be used",
	Long:  "set the project under which all applicable pctl command will be executed",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			err := fmt.Errorf("invalid command")
			exit.SetExitWithError(err, fmt.Sprintf("Failed to get the command args"))
			return
		}
		project := args[0]
		log.GetLogger().Debugf("Start [%s %s]", cmd.CommandPath(), project)
		_, err := testProjectAccess(project)
		if err != nil {
			exit.SetExitWithError(err, fmt.Sprintf("Failed to set project %s", project))
			return
		}
		conf := config.GetConfig()
		conf.Project = project
		configFile := context.GetContext().ConfigFilename()
		err = conf.Save(configFile)
		if err != nil {
			log.GetLogger().Debugf("Error in saving config file, [ERROR: %s]", err)
			exit.SetExitWithError(err, fmt.Sprintf("Error saving confing file"))
			return
		}

		log.GetLogger().Infof("Save config to file %s", configFile)
		log.GetLogger().Debugf("End [%s]", cmd.CommandPath())
	},
}

func init() {
	//******************
	// Config commands
	//
	//******************
	ConfigShowCmd.Flags().BoolP("source", "s", false, "Display the source of each configuration entry")
	ConfigCmd.AddCommand(
		ConfigInitCmd,
		ConfigShowCmd,
		ConfigSetCmd,
		// ConfigCheckCmd,
		newConfigDownloadCmd(commands.NewDownloadConfigOptions(log.GetLogger())),
	)
	ConfigSetCmd.AddCommand(ConfigSetProjectCmd)
}

func testProjectAccess(project string) (string, error) {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/project/%s/cluster", project)
	return auth.AuthAndRequest(uri, "GET", nil)
}
