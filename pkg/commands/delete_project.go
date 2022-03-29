package commands

import (
	"fmt"

	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/project"
	"github.com/RafayLabs/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteProjectOptions struct {
	YamlConfigPath string
	logger         log.Logger
	config         *config.Config
}

func NewDeleteProjectOptions(logger log.Logger, config *config.Config) *DeleteProjectOptions {
	o := new(DeleteProjectOptions)
	o.logger = logger
	o.config = config
	return o
}

func (o *DeleteProjectOptions) Validate(cmd *cobra.Command, args []string) error {
	flagSet := cmd.Flags()
	if flagSet.Changed("v3") {
		if flagSet.Changed(YamlConfigFlag) {
			if !utils.FileExists(o.YamlConfigPath) {
				return fmt.Errorf("file %s does not exist", o.YamlConfigPath)
			}
			return cobra.ExactArgs(0)(cmd, args)
		}
		return fmt.Errorf("invalid command, did not find flag '-f'")
	}
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	if ps, dup := utils.StringSliceContainsDuplicates(args); dup {
		return fmt.Errorf("project %s is given more than once", ps)
	}
	return nil
}

func (o *DeleteProjectOptions) Run(cmd *cobra.Command, args []string) error {
	logger := o.logger
	logger.Debugf("Start [%s]", cmd.CommandPath())
	// get the projects to delete
	for _, pr := range args {
		err := project.DeleteProject(pr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *DeleteProjectOptions) AddFlags(_ *cobra.Command) {}
