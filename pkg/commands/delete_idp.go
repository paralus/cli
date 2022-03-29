package commands

import (
	"fmt"

	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/idp"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteIdpOptions struct {
	logger log.Logger
	config *config.Config
}

func NewDeleteIdpOptions(logger log.Logger, config *config.Config) *DeleteIdpOptions {
	o := new(DeleteIdpOptions)
	o.logger = logger
	o.config = config
	return o
}

func (o *DeleteIdpOptions) Validate(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	if ps, dup := utils.StringSliceContainsDuplicates(args); dup {
		return fmt.Errorf("idp %s is given more than once", ps)
	}
	return nil
}

func (o *DeleteIdpOptions) Run(cmd *cobra.Command, args []string) error {
	logger := o.logger
	logger.Debugf("Start [%s]", cmd.CommandPath())

	// get the IDP to delete
	for _, i := range args {
		err := idp.DeleteIdp(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *DeleteIdpOptions) AddFlags(_ *cobra.Command) {}
