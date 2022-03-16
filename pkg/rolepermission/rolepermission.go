package rolepermission

import (
	"encoding/json"
	"fmt"

	"github.com/RafaySystems/rcloud-cli/pkg/config"
	"github.com/RafaySystems/rcloud-cli/pkg/models"
	"github.com/RafaySystems/rcloud-cli/pkg/rerror"
	"github.com/RafaySystems/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func ListRolePermissionWithCmd(cmd *cobra.Command) (*models.RolePermissionList, error) {
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
	rps := &models.RolePermissionList{}
	err = json.Unmarshal([]byte(resp), rps)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return rps, nil
}
