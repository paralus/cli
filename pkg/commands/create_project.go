package commands

import (
	"fmt"

	"github.com/RafaySystems/rcloud-cli/pkg/config"
	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"github.com/RafaySystems/rcloud-cli/pkg/project"
	"github.com/RafaySystems/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	CreateProjectDescriptionFlag = "desc"
)

type CreateProjectOptions struct {
	YamlConfigPath string
	Description    string
	logger         log.Logger
	config         *config.Config
}

func NewCreateProjectOptions(logger log.Logger, config *config.Config) *CreateProjectOptions {
	o := new(CreateProjectOptions)
	o.logger = logger
	o.config = config
	return o
}

func (o *CreateProjectOptions) Validate(cmd *cobra.Command, args []string) error {
	flagSet := cmd.Flags()
	if flagSet.Changed(YamlConfigFlag) {
		if !utils.FileExists(o.YamlConfigPath) {
			return fmt.Errorf("file %s does not exist", o.YamlConfigPath)
		}
		return cobra.ExactArgs(0)(cmd, args)
	}
	return cobra.ExactArgs(1)(cmd, args)
}

func (o *CreateProjectOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())
	projectName := args[0]
	return project.CreateProject(projectName, o.Description)

}

func (o *CreateProjectOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.PersistentFlags()
	flagSet.StringVar(&o.Description, CreateProjectDescriptionFlag, "",
		"Description for the project. Optional")

}
