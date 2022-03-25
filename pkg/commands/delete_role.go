package commands

import (
	"fmt"

	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/role"
	"github.com/RafayLabs/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteRoleOptions struct {
	logger log.Logger
	config *config.Config
}

func NewDeleteRoleOptions(logger log.Logger, config *config.Config) *DeleteRoleOptions {
	o := new(DeleteRoleOptions)
	o.logger = logger
	o.config = config
	return o
}

func (o *DeleteRoleOptions) Validate(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	if ps, dup := utils.StringSliceContainsDuplicates(args); dup {
		return fmt.Errorf("role %s is given more than once", ps)
	}
	return nil
}

func (o *DeleteRoleOptions) Run(cmd *cobra.Command, args []string) error {
	logger := o.logger
	logger.Debugf("Start [%s]", cmd.CommandPath())

	// get the role to delete
	for _, i := range args {
		err := role.DeleteRole(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *DeleteRoleOptions) AddFlags(_ *cobra.Command) {}
