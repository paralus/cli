package commands

import (
	"github.com/rafaylabs/rcloud-cli/pkg/cluster"
	"github.com/rafaylabs/rcloud-cli/pkg/config"
	"github.com/rafaylabs/rcloud-cli/pkg/log"
	"github.com/spf13/cobra"
)

type UpdateClusterOptions struct {
	YamlConfigPath string
	logger         log.Logger
}

func NewUpdateClusterOptions(logger log.Logger, config *config.Config) *UpdateClusterOptions {
	o := new(UpdateClusterOptions)
	o.logger = logger
	return o
}

func (o *UpdateClusterOptions) Validate(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}

	return nil
}

func (o *UpdateClusterOptions) Run(cmd *cobra.Command, args []string) error {
	updateCluster := false

	// retrieve the project id
	p, err := config.GetProjectIdFromFlagAndConfig(cmd)
	if err != nil {
		return err
	}

	// retrieve the original details
	c, err := cluster.GetCluster(args[0], p)
	if err != nil {
		return err
	}

	// get flagset
	if updateCluster {
		err = cluster.UpdateCluster(c)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *UpdateClusterOptions) AddFlags(_ *cobra.Command) {}
