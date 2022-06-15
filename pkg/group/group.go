package group

import (
	"encoding/json"
	"fmt"

	"github.com/oliveagle/jsonpath"
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/log"
	"github.com/paralus/cli/pkg/output"
	"github.com/paralus/cli/pkg/prefix"
	"github.com/paralus/cli/pkg/rerror"
	"github.com/paralus/cli/pkg/utils"
	commonv3 "github.com/paralus/paralus/proto/types/commonpb/v3"
	groupv3 "github.com/paralus/paralus/proto/types/userpb/v3"

	"github.com/spf13/cobra"
)

func ListGroupsWithCmd(cmd *cobra.Command) (*groupv3.GroupList, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/groups", cfg.Partner, cfg.Organization)
	uri = utils.AddPagenationToRequestWithCmd(cmd, uri)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, rerror.CrudErr{
			Type: "group",
			Name: "",
			Op:   "list",
		}
	}
	groups := &groupv3.GroupList{}
	err = json.Unmarshal([]byte(resp), groups)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return groups, nil
}

func GetGroupByName(groupName string) (*groupv3.Group, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/group/%s", cfg.Partner, cfg.Organization, groupName)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	grp := &groupv3.Group{}
	err = json.Unmarshal([]byte(resp), grp)
	if err != nil {
		return nil, err
	}

	return grp, nil
}

func getCount(obj interface{}) int {
	countInt := 0
	count, err := jsonpath.JsonPathLookup(obj, "$.count")
	if err == nil {
		if countF, ok := count.(float64); !ok {
			log.GetLogger().Infof("Failed to convert 'count' into float64")
			countInt = 0
		} else {
			countInt = int(countF)
		}
	} else {
		log.GetLogger().Debugf("Failed to get 'count' from result %s", err)
	}

	return countInt
}

func newGroupListSpec(obj interface{}) *output.OutputListSpec {
	countInt := getCount(obj)

	log.GetLogger().Debugf("count = %d", countInt)
	columnSpec := []output.ColumnSpec{
		{Header: "NAME", JsonPath: "name"},
	}

	spec := &output.OutputListSpec{
		NRows:   countInt,
		Base:    "$.results[%d]",
		Columns: columnSpec,
	}

	return spec
}

func NewGroupFromResponse(json_data []byte) (*groupv3.Group, error) {
	var gr groupv3.Group
	if err := json.Unmarshal(json_data, &gr); err != nil {
		return nil, err
	}
	return &gr, nil
}

func ListPrint(cmd *cobra.Command, jsonObj []byte) error {
	return output.PrintList(cmd, jsonObj, newGroupListSpec)
}

func Print(cmd *cobra.Command, jsonObj []byte) error {
	g, err := NewGroupFromResponse(jsonObj)
	if err != nil {
		return err
	}

	w := prefix.NewPrefixWriter()
	w.Write(0, "Name: %s\n", g.Metadata.Name)
	w.Write(0, "Description: %s\n", g.Metadata.Description)
	w.Write(0, "Type: %s\n", g.Spec.Type)
	w.Write(0, "Users: %s\n", g.Spec.Users)

	w.Flush()

	return nil
}

func CreateGroup(name, description string) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	group := groupv3.Group{
		Kind: "Group",
		Metadata: &commonv3.Metadata{
			Name:        name,
			Description: description,
		},
	}
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/groups", cfg.Partner, cfg.Organization)
	_, err := auth.AuthAndRequest(uri, "POST", group)
	if err != nil {
		return err
	}
	return nil
}

func DeleteGroup(groupName string) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/group/%s", cfg.Partner, cfg.Organization, groupName)
	_, err := auth.AuthAndRequest(uri, "DELETE", nil)
	return err
}

// Update group takes the updated group details and sends it to the core
func UpdateGroup(grp *groupv3.Group) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/group/%s", cfg.Partner, cfg.Organization, grp.Metadata.Name)
	_, err := auth.AuthAndRequest(uri, "PUT", grp)
	if err != nil {
		return rerror.CrudErr{
			Type: "group",
			Name: grp.Metadata.Name,
			Op:   "update",
		}
	}
	return nil
}

// Apply group takes the group details and sends it to the core
func ApplyGroup(grp *groupv3.Group) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()

	grpExisting, _ := GetGroupByName(grp.Metadata.Name)
	if grpExisting != nil {
		log.GetLogger().Debugf("updating group: %s", grp.Metadata.Name)
		uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/group/%s", cfg.Partner, cfg.Organization, grp.Metadata.Name)
		_, err := auth.AuthAndRequest(uri, "PUT", grp)
		if err != nil {
			return rerror.CrudErr{
				Type: "group",
				Name: grp.Metadata.Name,
				Op:   "update",
			}
		}
	} else {
		log.GetLogger().Debugf("creating group: %s", grp.Metadata.Name)
		uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/groups", cfg.Partner, cfg.Organization)
		_, err := auth.AuthAndRequest(uri, "POST", grp)
		if err != nil {
			return rerror.CrudErr{
				Type: "group",
				Name: grp.Metadata.Name,
				Op:   "create",
			}
		}
	}
	return nil
}
