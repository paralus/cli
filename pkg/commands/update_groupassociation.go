package commands

import (
	"fmt"
	"regexp"

	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/group"
	"github.com/paralus/cli/pkg/log"
	"github.com/paralus/cli/pkg/project"
	"github.com/paralus/cli/pkg/role"
	"github.com/paralus/cli/pkg/user"
	userv3 "github.com/paralus/paralus/proto/types/userpb/v3"
	"github.com/spf13/cobra"
)

const (
	UpdateGroupProjectFlag    = "associateproject"
	UpdateGroupRolesFlag      = "roles"
	UpdateNamespacesFlag      = "namespace"
	UpdateGroupUserFlag       = "associateuser"
	UpdateAddGroupUserFlag    = "addusers"
	UpdateRemoveGroupUserFlag = "removeusers"
)

type UpdateGroupassociationOptions struct {
	Project     string
	Roles       []string
	Scope       string
	Namespace   string
	User        string
	AddUsers    []string
	RemoveUsers []string
	logger      log.Logger
	config      *config.Config
}

func NewUpdateGroupassociationOptions(logger log.Logger, config *config.Config) *UpdateGroupassociationOptions {
	o := new(UpdateGroupassociationOptions)
	o.logger = logger
	o.config = config
	return o
}
func (o *UpdateGroupassociationOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MinimumNArgs(1)(cmd, args)
}

func (o *UpdateGroupassociationOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())
	flagSet := cmd.Flags()
	name := args[0]
	err := fmt.Errorf("flags not triggered")
	if flagSet.Changed(UpdateGroupUserFlag) {
		err = UpdateUserAssociation(cmd, name, o.AddUsers, o.RemoveUsers)
		if err != nil {
			return err
		}
	} else if flagSet.Changed(UpdateGroupProjectFlag) {
		err = UpdateProjectAssociation(cmd, name, o.Project, o.Roles, o.Namespace)
		if err != nil {
			return err
		}
	}

	return err
}

func UpdateProjectAssociation(cmd *cobra.Command, groupName string, projectName string, chosenRoles []string, namespace string) error {
	currGroup, err := group.GetGroupByName(groupName)
	if err != nil {
		return fmt.Errorf("group %s does not exist", groupName)
	}

	projectResp, err := project.GetProjectByName(projectName)
	if err != nil {
		return fmt.Errorf("project %s does not exist", projectName)
	}

	roleList, err := role.ListRolesWithCmd(cmd)
	if err != nil {
		return fmt.Errorf("failed to fetch role details, cause: %s", err.Error())
	}

	if len(currGroup.Spec.ProjectNamespaceRoles) == 0 {
		currGroup.Spec.ProjectNamespaceRoles = make([]*userv3.ProjectNamespaceRole, 0)
	}

	regexc := regexp.MustCompile(`[^a-z0-9-]+`)

	for _, eachRole := range roleList.Items {
		if StringInSlice(eachRole.Metadata.Name, chosenRoles) {
			if eachRole.Spec.Scope == "namespace" {
				if namespace == "" {
					return fmt.Errorf("namespace not specified for a namespaced role")
				}
				match := regexc.MatchString(namespace)
				if match {
					return fmt.Errorf("namespace %q is invalid", namespace)
				}
				if !(len(namespace) >= 1 && len(namespace) <= 63) {
					return fmt.Errorf("namespace %q is invalid. must be no more than 63 characters", namespace)
				}
				currGroup.Spec.ProjectNamespaceRoles = append(currGroup.Spec.ProjectNamespaceRoles, &userv3.ProjectNamespaceRole{
					Project:   &projectResp.Metadata.Name,
					Role:      eachRole.Metadata.Name,
					Namespace: &namespace,
				})
			} else {
				currGroup.Spec.ProjectNamespaceRoles = append(currGroup.Spec.ProjectNamespaceRoles, &userv3.ProjectNamespaceRole{
					Project: &projectResp.Metadata.Name,
					Role:    eachRole.Metadata.Name,
				})
			}
		}
	}

	if len(chosenRoles) == 0 {
		currGroup.Spec.ProjectNamespaceRoles = append(currGroup.Spec.ProjectNamespaceRoles, &userv3.ProjectNamespaceRole{
			Project: &projectResp.Metadata.Name,
		})
	}

	err = group.UpdateGroup(currGroup)
	if err != nil {
		return fmt.Errorf("unable to update group %s, cause: %s", currGroup.Metadata.Name, err.Error())
	}
	return err
}

func UpdateUserAssociation(cmd *cobra.Command, groupName string, addUsernames []string, removeUsernames []string) error {
	currGroup, err := group.GetGroupByName(groupName)
	if err != nil {
		return fmt.Errorf("group %s does not exist", groupName)
	}

	for _, eachAddUser := range addUsernames {
		usr, err := user.GetUserByName(eachAddUser)
		if err != nil {
			return fmt.Errorf("user %s does not exist", eachAddUser)
		}
		currGroup.Spec.Users = append(currGroup.Spec.Users, usr.Metadata.Name)
	}

	for _, eachRemoveUser := range removeUsernames {
		usr, err := user.GetUserByName(eachRemoveUser)
		if err != nil {
			return fmt.Errorf("user %s does not exist", eachRemoveUser)
		}

		findAndDelete(currGroup.Spec.Users, usr.Metadata.Name)
	}

	err = group.UpdateGroup(currGroup)
	if err != nil {
		return fmt.Errorf("unable to update group %s, cause: %s", currGroup.Metadata.Name, err.Error())
	}

	return nil

}

func findAndDelete(s []string, item string) []string {
	index := 0
	for _, i := range s {
		if i != item {
			s[index] = i
			index++
		}
	}
	return s[:index]
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (o *UpdateGroupassociationOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.PersistentFlags()
	flagSet.StringVar(&o.Project, UpdateGroupProjectFlag, "",
		"Project to assign to the Group.")
	flagSet.StringSliceVar(&o.Roles, UpdateGroupRolesFlag, nil,
		"Select Roles to assign for Project.")
	flagSet.StringVar(&o.Namespace, UpdateNamespacesFlag, "",
		"Select Namespace to assign for Project.")
	//User Flags
	flagSet.StringVar(&o.User, UpdateGroupUserFlag, "",
		"Declare associated users to be edited.")
	flagSet.StringSliceVar(&o.AddUsers, UpdateAddGroupUserFlag, nil,
		"Select Users to assign for the group. Optional")
	flagSet.StringSliceVar(&o.RemoveUsers, UpdateRemoveGroupUserFlag, nil,
		"Select Users to remove from the group. Optional")
}
