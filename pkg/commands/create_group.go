package commands

import (
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/group"
	"github.com/paralus/cli/pkg/log"
	"github.com/spf13/cobra"
)

const (
	CreateGroupDescriptionFlag = "desc"
)

type CreateGroupOptions struct {
	Description string
	logger      log.Logger
	config      *config.Config
}

func NewCreateGroupOptions(logger log.Logger, config *config.Config) *CreateGroupOptions {
	o := new(CreateGroupOptions)
	o.logger = logger
	o.config = config
	return o
}

func (o *CreateGroupOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.ExactArgs(1)(cmd, args)
}

func (o *CreateGroupOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())

	GroupName := args[0]

	err := group.CreateGroup(GroupName, o.Description)
	if err != nil {
		return err
	}

	return err
}

func (o *CreateGroupOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.PersistentFlags()
	flagSet.StringVar(&o.Description, CreateGroupDescriptionFlag, "",
		"Description for the Group. Optional")
}
