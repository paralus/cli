package commands

import (
	"fmt"

	"github.com/paralus/cli/pkg/log"
	"github.com/paralus/cli/pkg/user"
	"github.com/spf13/cobra"
)

type DeleteUserOptions struct {
	logger log.Logger
}

func NewDeleteUserOptions(logger log.Logger) *DeleteUserOptions {
	o := new(DeleteUserOptions)
	o.logger = logger
	return o
}
func (o *DeleteUserOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.ExactArgs(1)(cmd, args)
}

func (o *DeleteUserOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())

	username := args[0]

	err := DeleteUser(cmd, username)
	if err != nil {
		return err
	}

	return err
}

func DeleteUser(cmd *cobra.Command, username string) error {
	_, err := user.GetUserByName(username)
	if err != nil {
		return fmt.Errorf("username input is invalid")
	}

	err = user.DeleteUser(username)
	if err != nil {
		return fmt.Errorf("unable to delete user %s, cause: %s", username, err.Error())
	}
	return err
}

func (c *DeleteUserOptions) AddFlags(_ *cobra.Command) {}
