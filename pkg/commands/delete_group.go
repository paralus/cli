package commands

import (
	"fmt"

	"github.com/RafaySystems/rcloud-cli/pkg/config"
	"github.com/RafaySystems/rcloud-cli/pkg/group"
	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"github.com/RafaySystems/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteGroupOptions struct {
	logger log.Logger
	config *config.Config
}

func NewDeleteGroupOptions(logger log.Logger, config *config.Config) *DeleteGroupOptions {
	o := new(DeleteGroupOptions)
	o.logger = logger
	o.config = config
	return o
}

func (o *DeleteGroupOptions) Validate(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	if ps, dup := utils.StringSliceContainsDuplicates(args); dup {
		return fmt.Errorf("group %s is given more than once", ps)
	}
	return nil
}

func (o *DeleteGroupOptions) Run(cmd *cobra.Command, args []string) error {
	logger := o.logger
	logger.Debugf("Start [%s]", cmd.CommandPath())

	// get the groups to delete
	for _, gr := range args {
		err := group.DeleteGroup(gr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *DeleteGroupOptions) AddFlags(_ *cobra.Command) {}
