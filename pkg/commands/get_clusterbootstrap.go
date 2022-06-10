package commands

import (
	"fmt"

	"github.com/paralus/cli/pkg/cluster"
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/log"
	"github.com/spf13/cobra"
)

type GetClusterBootstrapOptions struct {
	logger log.Logger
}

func NewGetClusterBootstrapOptions(logger log.Logger) *GetClusterBootstrapOptions {
	o := new(GetClusterBootstrapOptions)
	o.logger = logger
	return o
}

func (o *GetClusterBootstrapOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.ExactArgs(1)(cmd, args)
}

func (o *GetClusterBootstrapOptions) Run(cmd *cobra.Command, args []string) error {
	// retrieve the project
	p, err := config.GetProjectIdFromFlagAndConfig(cmd)
	if err != nil {
		return err
	}

	// get cluster name
	n := args[0]

	// get the bootstrap file
	bf, err := cluster.GetBootstrapFile(n, p)
	if err != nil {
		return err
	}
	fmt.Println(bf)
	return nil
}

func (o *GetClusterBootstrapOptions) AddFlags(_ *cobra.Command) {}
