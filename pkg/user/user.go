package user

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
	userv3 "github.com/paralus/paralus/proto/types/userpb/v3"
	"github.com/spf13/cobra"
)

func ListUsersWithCmd(cmd *cobra.Command) (*userv3.UserList, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := "/auth/v3/users"
	uri = utils.AddPagenationToRequestWithCmd(cmd, uri)
	uri = uri + fmt.Sprintf("&partner=%s&organization=%s", cfg.Partner, cfg.Organization)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, rerror.CrudErr{
			Type: "user",
			Name: "",
			Op:   "list",
		}
	}
	users := &userv3.UserList{}
	err = json.Unmarshal([]byte(resp), users)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return users, nil
}

func GetUserByName(userName string) (*userv3.User, error) {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/user/%s", userName)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	user := &userv3.User{}
	err = json.Unmarshal([]byte(resp), user)
	if err != nil {
		return nil, err
	}

	return user, nil

}

func NewUserFromResponse(json_data []byte) (*userv3.User, error) {
	var ur userv3.User
	if err := json.Unmarshal(json_data, &ur); err != nil {
		return nil, err
	}
	fmt.Println(ur)
	return &ur, nil
}

func CreateUser(usr *userv3.User) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := "/auth/v3/users"
	_, err := auth.AuthAndRequest(uri, "POST", usr)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(userName string) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/user/%s", userName)
	_, err := auth.AuthAndRequest(uri, "DELETE", nil)
	return err
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

func newUserListSpec(obj interface{}) *output.OutputListSpec {
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

func ListPrint(cmd *cobra.Command, jsonObj []byte) error {
	return output.PrintList(cmd, jsonObj, newUserListSpec)
}

func Print(cmd *cobra.Command, jsonObj []byte) error {
	g, err := NewUserFromResponse(jsonObj)
	if err != nil {
		return err
	}

	w := prefix.NewPrefixWriter()
	w.Write(0, "Username: %s\n", g.Metadata.Name)
	w.Write(0, "First Name: %s\n", g.Spec.FirstName)
	w.Write(0, "Last Name: %s\n", g.Spec.LastName)
	w.Write(0, "Groups: %s\n", g.Spec.Groups)

	w.Flush()

	return nil
}

// Apply user takes the user details and sends it to the core
func ApplyUser(usr *userv3.User) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()

	projExisting, _ := GetUserByName(usr.Metadata.Name)
	if projExisting != nil {
		log.GetLogger().Debugf("updating user: %s", usr.Metadata.Name)
		uri := fmt.Sprintf("/auth/v3/user/%s", usr.Metadata.Name)
		_, err := auth.AuthAndRequest(uri, "PUT", usr)
		if err != nil {
			return rerror.CrudErr{
				Type: "user",
				Name: usr.Metadata.Name,
				Op:   "update",
			}
		}
	} else {
		log.GetLogger().Debugf("creating user: %s", usr.Metadata.Name)
		uri := "/auth/v3/users/"
		_, err := auth.AuthAndRequest(uri, "POST", usr)
		if err != nil {
			return rerror.CrudErr{
				Type: "user",
				Name: usr.Metadata.Name,
				Op:   "create",
			}
		}
	}
	return nil
}
