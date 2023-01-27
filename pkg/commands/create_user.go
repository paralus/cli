package commands

import (
	"fmt"

	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/group"
	"github.com/paralus/cli/pkg/log"
	"github.com/paralus/cli/pkg/user"
	commonv3 "github.com/paralus/paralus/proto/types/commonpb/v3"
	userv3 "github.com/paralus/paralus/proto/types/userpb/v3"
	"github.com/spf13/cobra"
)

const (
	CreateConsoleAccessFlag  = "console"
	CreateUserGroupAssocFlag = "groups"
)

type CreateUserOptions struct {
	ConsoleAccessInputs []string
	Groups              []string
	logger              log.Logger
	config              *config.Config
}

func NewCreateUserOptions(logger log.Logger, config *config.Config) *CreateUserOptions {
	o := new(CreateUserOptions)
	o.logger = logger
	o.config = config
	return o
}
func (o *CreateUserOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.ExactArgs(1)(cmd, args)
}

func (o *CreateUserOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())

	Username := args[0]

	err := CreateUser(cmd, Username, o.Groups, o.ConsoleAccessInputs)
	if err != nil {
		return err
	}
	return nil
}

func CreateUser(cmd *cobra.Command, username string, groups []string, ConsoleAccessInputs []string) error {
	flags := cmd.Flags()

	newAccount := userv3.User{
		Kind: "User",
		Metadata: &commonv3.Metadata{
			Name:         username,
			Organization: config.GetConfig().Organization,
			Partner:      config.GetConfig().Partner,
		},
		Spec: &userv3.UserSpec{
			ForceReset: true,
		},
	}
	if flags.Changed(CreateConsoleAccessFlag) {
		if len(ConsoleAccessInputs) <= 2 {
			newAccount.Spec = &userv3.UserSpec{
				FirstName: ConsoleAccessInputs[0],
				LastName:  ConsoleAccessInputs[1],
				Phone:     "",
			}
		} else {
			newAccount.Spec = &userv3.UserSpec{
				FirstName: ConsoleAccessInputs[0],
				LastName:  ConsoleAccessInputs[1],
				Phone:     ConsoleAccessInputs[2],
			}
		}
	}

	defaultAssignedGroup := "All Local Users"
	grp, err := group.GetGroupByName(defaultAssignedGroup)
	if err != nil {
		return fmt.Errorf("group %s does not exist", defaultAssignedGroup)
	}
	newAccount.Spec.Groups = append(newAccount.Spec.Groups, grp.Metadata.Name)

	if flags.Changed(CreateUserGroupAssocFlag) {
		for _, groupName := range groups {
			grp, err = group.GetGroupByName(groupName)
			if err != nil {
				return fmt.Errorf("group %s does not exist", groupName)
			}
			newAccount.Spec.Groups = append(newAccount.Spec.Groups, grp.Metadata.Name)

		}
	}

	err = user.CreateUser(&newAccount)
	if err != nil {
		return fmt.Errorf("user creation failed, cause: %s", err.Error())
	}
	return err

}

func (o *CreateUserOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.PersistentFlags()
	flagSet.StringSliceVar(&o.Groups, CreateUserGroupAssocFlag, nil,
		"Existing Groups to be assigned to. Optional")
	flagSet.StringSliceVar(&o.ConsoleAccessInputs, CreateConsoleAccessFlag, nil,
		"Console access flag. Optional, but if Console access flag is enabled, firstname and lastname inputs are required, while phone number input is optional.")
}
