package role

import (
	"encoding/json"
	"fmt"

	"github.com/RafaySystems/rcloud-cli/pkg/config"
	"github.com/RafaySystems/rcloud-cli/pkg/models"
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
