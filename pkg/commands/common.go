package commands

import (
	"fmt"
	"path/filepath"

	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/constants"
	"github.com/RafayLabs/rcloud-cli/pkg/context"
	"github.com/RafayLabs/rcloud-cli/pkg/exit"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/output"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	YamlConfigFlag          = "config-file"
	YamlConfigShorthandFlag = "f"
)

// Supported Resources
var resources = []string{"Cluster", "Project", "User", "Group", "Idp", "OIDCProvider", "Role"}

func AddYamlConfigFlag(cmd *cobra.Command, flagHelp string) {
	cmd.PersistentFlags().StringP(YamlConfigFlag, YamlConfigShorthandFlag, "", flagHelp)
	cmd.MarkPersistentFlagFilename(YamlConfigFlag, "yml", "yaml")
}

func AddYamlConfigFlagVar(ptr *string, cmd *cobra.Command, flagHelp string) {
	cmd.PersistentFlags().StringVarP(ptr, YamlConfigFlag, YamlConfigShorthandFlag, "", flagHelp)
	cmd.MarkPersistentFlagFilename(YamlConfigFlag, "yml", "yaml")
}

type CmdOptions interface {
	// Validate is used to validate arguments and flags.
	// The function will validate args without calling REST APIs.
	// This means validating if files exist, if there are duplicate arguments provided,
	// if the proper flags are provided, etc.
	// It is also where config files are parsed
	Validate(cmd *cobra.Command, args []string) error
	// Run runs the command action
	Run(cmd *cobra.Command, args []string) error
	// AddFlags adds flags to the supplied cobra command
	AddFlags(cmd *cobra.Command)
}

// GlobalOptions is a struct to hold the values for global options
type GlobalOptions struct {
	Verbose,
	Debug bool
	ConfigFile,
	Output,
	Project string
	FilePath string
	Group    string
	wait     bool
	config   *config.Config
	logger   log.Logger
}

type gvkYamlConfig struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

func NewGlobalOptions(log log.Logger, config *config.Config) *GlobalOptions {
	globalOptions := new(GlobalOptions)
	globalOptions.config = config
	globalOptions.logger = log
	return globalOptions
}

func (g *GlobalOptions) Validate(_ *cobra.Command, _ []string) error {
	return nil
}

func (g *GlobalOptions) Run(cmd *cobra.Command, _ []string) error {
	isVerbose, err := cmd.Flags().GetBool(constants.VERBOSE_FLAG_NAME)
	if err != nil {
		isVerbose = false
	}

	isDebug, err := cmd.Flags().GetBool(constants.DEBUG_FLAG_NAME)
	if err != nil {
		isDebug = false
	}

	// set the desired logging level
	// by default, the log level is set to error
	if isVerbose {
		log.SetLevel(zap.InfoLevel)
	}
	if isDebug {
		log.SetLevel(zap.DebugLevel)
	}

	isStructuredOutput, err := cmd.Flags().GetBool(constants.STRUCTURED_OUTPUT_FLAG_NAME)
	if err != nil {
		isStructuredOutput = false
	}

	configFile, err := cmd.Flags().GetString(constants.CONFIG_FLAG_NAME)
	if err != nil {
		return err
	}

	cliCtx := context.GetContext()
	if len(configFile) != 0 {
		log.GetLogger().Infof("Config options is set up %s \n", configFile)

		dir, file := filepath.Split(configFile)
		if len(dir) != 0 {
			cliCtx.ConfigDir = dir
		}
		if len(file) != 0 {
			cliCtx.ConfigFile = file
		}
	}

	cliCtx.Verbose = isVerbose
	cliCtx.Debug = isDebug
	cliCtx.StructuredOutput = isStructuredOutput

	log.GetLogger().Debugf("Prerun")

	// Log the context
	cliCtx.Log("Context: ")
	err = config.InitConfig(cliCtx)
	if err != nil {
		path := cliCtx.ConfigFilename()

		if cmd.CommandPath() != "rctl config init" {
			// Ignore this error in case of 'config init', this error will be ignored.
			exit.SetExitWithError(err, fmt.Sprintf("Failed to load config file from %s. "+
				"Please use 'rctl config init' to install the config file", path))
			output.Exit()
		}
	}
	if cmd.CommandPath() != "rctl config init" && cmd.CommandPath() != "rctl version" {
		err := g.config.MiniCheck()
		if err != nil {
			// Ignore this error in case of 'config init' or 'version', this error will be ignored.
			exit.SetExitWithError(err, "Config check failed")
			output.Exit()
		}
	}

	// Log the config
	g.config.Log("Config: ")

	// check if the project flag is provided
	if cmd.Flags().Changed("project") {
		_, err := config.GetProjectIdByName(g.Project)
		if err != nil {
			return fmt.Errorf("invalid flag/argument passed for project")
		}
		g.config.Project = g.Project
	}

	// check if the wait flag is provided
	if cmd.Flags().Changed("wait") {
		g.wait = true
	}
	return nil
}

func (g *GlobalOptions) AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&g.Verbose, "verbose", "v", false, "Verbose mode. A lot more information output.")
	cmd.PersistentFlags().BoolVarP(&g.Debug, "debug", "d", false, "Enable debug logs")
	cmd.PersistentFlags().BoolVarP(&g.wait, "wait", "", false, "Wait for cluster info to publish namespace")
	cmd.PersistentFlags().StringVarP(&g.ConfigFile, "config", "c", "", "Customize cli config file")
	cmd.PersistentFlags().StringVarP(&g.Output, "output", "o", "table", "Print json, yaml or table output. Default is table")
	cmd.PersistentFlags().StringVarP(&g.Project, "project", "p", "", "provide a specific project context")
	cmd.PersistentFlags().StringVarP(&g.FilePath, "file", "f", "", "provide file with resource to be created")
}

func getProjectIdFromFlag(cmd *cobra.Command) (string, error) {
	projectName, err := cmd.Flags().GetString("project")
	if err != nil {
		exit.SetExitWithError(err, "Invalid flag/argument passed for project")
		return "", err
	}
	if projectName != "" {
		return config.GetProjectIdByName(projectName)
	}
	return "", nil
}
