package role

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/RafaySystems/rcloud-cli/pkg/config"
	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"github.com/RafaySystems/rcloud-cli/pkg/models"
	"github.com/RafaySystems/rcloud-cli/pkg/prefix"
	"github.com/RafaySystems/rcloud-cli/pkg/rerror"
	"github.com/RafaySystems/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func ListRolesWithCmd(cmd *cobra.Command) (*models.RoleList, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/roles", cfg.Partner, cfg.Organization)
	uri = utils.AddPagenationToRequestWithCmd(cmd, uri)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, rerror.CrudErr{
			Type: "role",
			Name: "",
			Op:   "list",
		}
	}
	rls := &models.RoleList{}
	err = json.Unmarshal([]byte(resp), rls)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return rls, nil
}

func CreateRole(r *models.Role) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	//set partner and organization
	r.Metadata.Partner = cfg.Partner
	r.Metadata.Organization = cfg.Organization

	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/roles", cfg.Partner, cfg.Organization)
	_, err := auth.AuthAndRequest(uri, "POST", r)
	if err != nil {
		return err
	}
	return nil
}

func DeleteRole(name string) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/role/%s", cfg.Partner, cfg.Organization, name)
	_, err := auth.AuthAndRequest(uri, "DELETE", nil)
	return err
}

func GetRoleByName(name string) (*models.Role, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/role/%s", cfg.Partner, cfg.Organization, name)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	r := &models.Role{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		return nil, err
	}

	return r, nil

}

// Apply role takes the role details and sends it to the core
func ApplyRole(r *models.Role) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	er, _ := GetRoleByName(r.Metadata.Name)
	if er != nil {
		log.GetLogger().Debugf("updating role: %s", r.Metadata.Name)
		//set partner and organization
		r.Metadata.Partner = cfg.Partner
		r.Metadata.Organization = cfg.Organization
		uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/role/%s", cfg.Partner, cfg.Organization, r.Metadata.Name)
		_, err := auth.AuthAndRequest(uri, "PUT", r)
		if err != nil {
			return rerror.CrudErr{
				Type: "role",
				Name: r.Metadata.Name,
				Op:   "update",
			}
		}
	} else {
		log.GetLogger().Debugf("creating role: %s", r.Metadata.Name)
		//set partner and organization
		r.Metadata.Partner = cfg.Partner
		r.Metadata.Organization = cfg.Organization
		uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/roles", cfg.Partner, cfg.Organization)
		_, err := auth.AuthAndRequest(uri, "POST", r)
		if err != nil {
			return rerror.CrudErr{
				Type: "role",
				Name: r.Metadata.Name,
				Op:   "create",
			}
		}
	}
	return nil
}

func NewRoleFromResponse(json_data []byte) (*models.Role, error) {
	var gr models.Role
	if err := json.Unmarshal(json_data, &gr); err != nil {
		return nil, err
	}
	return &gr, nil
}

func Print(cmd *cobra.Command, jsonObj []byte) error {
	r, err := NewRoleFromResponse(jsonObj)
	if err != nil {
		return err
	}

	w := prefix.NewPrefixWriter()
	w.Write(0, "Name: %s\n", r.Metadata.Name)
	w.Write(0, "Description: %s\n", r.Metadata.Description)
	w.Write(0, "Scope: %s\n", r.Spec.Scope)
	w.Write(0, "IsGlobal: %s\n", strconv.FormatBool(r.Spec.IsGlobal))
	w.Write(0, "Permissions: %s\n", r.Spec.Rolepermissions)

	w.Flush()

	return nil
}
