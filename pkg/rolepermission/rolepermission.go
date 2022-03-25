package rolepermission

import (
	"encoding/json"
	"fmt"

	rolev3 "github.com/RafayLabs/rcloud-base/proto/types/rolepb/v3"
	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/rerror"
	"github.com/RafayLabs/rcloud-cli/pkg/utils"
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
