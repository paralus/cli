package commands

import (
	"fmt"

	"github.com/rafaylabs/rcloud-cli/pkg/config"
	"github.com/rafaylabs/rcloud-cli/pkg/log"
	"github.com/rafaylabs/rcloud-cli/pkg/models"
	"github.com/rafaylabs/rcloud-cli/pkg/role"
	"github.com/rafaylabs/rcloud-cli/pkg/rolepermission"
	"github.com/spf13/cobra"
)

const (
	ScopeFlag       = "scope"
	PermissionsFlag = "permissions"
)

type CreateRoleOptions struct {
	scope       string
	permissions []string
	logger      log.Logger
	config      *config.Config
}

func NewCreateRoleOptions(logger log.Logger, config *config.Config) *CreateRoleOptions {
	o := new(CreateRoleOptions)
	o.logger = logger
	o.config = config
	return o
}

func (c *CreateRoleOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.PersistentFlags()
	flagSet.StringVar(&c.scope, ScopeFlag, "",
		"Scope of this role, can be either ORGANIZATION, PROJECT.")
	flagSet.StringSliceVar(&c.permissions, PermissionsFlag, nil,
		"Permissions that are to be associated to this role.")
}

func (o *CreateRoleOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MinimumNArgs(1)(cmd, args)

}

func (o *CreateRoleOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())

	name := args[0]
	flagSet := cmd.Flags()
	err := fmt.Errorf("flags not triggered")

	if flagSet.Changed(ScopeFlag) && flagSet.Changed(PermissionsFlag) {

		// get scope
		scope, err := flagSet.GetString(ScopeFlag)
		if err != nil {
			return err
		} else {
			if scope != "ORGANIZATION" && scope != "PROJECT" && scope != "" {
				return fmt.Errorf("scope can be either ORGANIZATION or PROJECT, given scope is %s ", scope)
			}
		}

		//validate provided permissions
		rps, err := rolepermission.ListRolePermissionWithCmd(cmd)
		if err != nil {
			return fmt.Errorf("unable to verify permissions, error: %s ", err.Error())
		} else {
			pl := []string{}
			for _, rp := range rps.Items {
				pl = append(pl, rp.Metadata.Name)
			}
			// check for invalid permissions
			for _, p := range o.permissions {
				if !contains(pl, p) {
					return fmt.Errorf("invalid permission %s ", p)
				}
			}
		}

		cr := &models.Role{
			Kind: "Role",
			Metadata: models.Metadata{
				Name: name,
			},
			Spec: models.RoleSpec{
				Rolepermissions: o.permissions,
				IsGlobal:        false,
				Scope:           scope,
			},
		}

		err = role.CreateRole(cr)
		if err != nil {
			return fmt.Errorf("failed to create role, error: %s", err.Error())
		}
	} else {
		return err
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, perm := range slice {
		if perm == item {
			return true
		}
	}
	return false
}
