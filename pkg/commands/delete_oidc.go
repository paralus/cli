package commands

import (
	"fmt"

	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/oidc"
	"github.com/RafayLabs/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteOIDCProviderOptions struct {
	logger log.Logger
	config *config.Config
}

func NewDeleteOIDCProviderOptions(logger log.Logger, config *config.Config) *DeleteOIDCProviderOptions {
	o := new(DeleteOIDCProviderOptions)
	o.logger = logger
	o.config = config
	return o
}

func (o *DeleteOIDCProviderOptions) Validate(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	if ps, dup := utils.StringSliceContainsDuplicates(args); dup {
		return fmt.Errorf("oidc %s is given more than once", ps)
	}
	return nil
}

func (o *DeleteOIDCProviderOptions) Run(cmd *cobra.Command, args []string) error {
	logger := o.logger
	logger.Debugf("Start [%s]", cmd.CommandPath())

	// get the oidc to delete
	for _, i := range args {
		err := oidc.DeleteOIDCProvider(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *DeleteOIDCProviderOptions) AddFlags(_ *cobra.Command) {}
