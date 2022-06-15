package rolepermission

import (
	"encoding/json"
	"fmt"

	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/rerror"
	"github.com/paralus/cli/pkg/utils"
	rolev3 "github.com/paralus/paralus/proto/types/rolepb/v3"
	"github.com/spf13/cobra"
)

func ListRolePermissionWithCmd(cmd *cobra.Command) (*rolev3.RolePermissionList, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := "/auth/v3/rolepermissions"
	uri = utils.AddPagenationToRequestWithCmd(cmd, uri)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, rerror.CrudErr{
			Type: "rolepermission",
			Name: "",
			Op:   "list",
		}
	}
	rps := &rolev3.RolePermissionList{}
	err = json.Unmarshal([]byte(resp), rps)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return rps, nil
}

func ListRolePermissionWithScope(cmd *cobra.Command, scope string) (*rolev3.RolePermissionList, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := "/auth/v3/rolepermissions" + fmt.Sprintf("?selector=%s", scope)
	uri = utils.AddPagenationToRequestWithCmd(cmd, uri)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, rerror.CrudErr{
			Type: "rolepermission",
			Name: "",
			Op:   "list",
		}
	}
	rps := &rolev3.RolePermissionList{}
	err = json.Unmarshal([]byte(resp), rps)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return rps, nil
}
