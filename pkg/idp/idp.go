package idp

import (
	"encoding/json"
	"fmt"

	systemv3 "github.com/RafayLabs/rcloud-base/proto/types/systempb/v3"
	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/output"
	"github.com/RafayLabs/rcloud-cli/pkg/prefix"
	"github.com/RafayLabs/rcloud-cli/pkg/rerror"
	"github.com/RafayLabs/rcloud-cli/pkg/utils"
	"github.com/oliveagle/jsonpath"
	"github.com/spf13/cobra"
)

func ListIdpWithCmd(cmd *cobra.Command) (*systemv3.IdpList, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := "/auth/v3/sso/idp"
	uri = utils.AddPagenationToRequestWithCmd(cmd, uri)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, rerror.CrudErr{
			Type: "user",
			Name: "",
			Op:   "list",
		}
	}
	idps := &systemv3.IdpList{}
	err = json.Unmarshal([]byte(resp), idps)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return idps, nil
}

func GetIdpByName(idpName string) (*systemv3.Idp, error) {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/sso/idp/%s", idpName)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	idp := &systemv3.Idp{}
	err = json.Unmarshal([]byte(resp), idp)
	if err != nil {
		return nil, err
	}

	return idp, nil

}

func NewIdpFromResponse(json_data []byte) (*systemv3.Idp, error) {
	var ur systemv3.Idp
	if err := json.Unmarshal(json_data, &ur); err != nil {
		return nil, err
	}
	fmt.Println(ur)
	return &ur, nil
}

func CreateIdp(idp *systemv3.Idp) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := "/auth/v3/sso/idp"
	_, err := auth.AuthAndRequest(uri, "POST", idp)
	if err != nil {
		return err
	}
	return nil
}

func DeleteIdp(idpName string) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/sso/idp/%s", idpName)
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
	g, err := NewIdpFromResponse(jsonObj)
	if err != nil {
		return err
	}

	w := prefix.NewPrefixWriter()
	w.Write(0, "Name: %s\n", g.Metadata.Name)
	w.Write(0, "Domain: %s\n", g.Spec.Domain)
	w.Write(0, "AcsURL: %s\n", g.Spec.AcsUrl)

	w.Flush()

	return nil
}

// Apply idp takes the idp details and sends it to the core
func ApplyIDP(id *systemv3.Idp) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()

	idpExisting, _ := GetIdpByName(id.Metadata.Name)
	if idpExisting != nil {
		log.GetLogger().Debugf("updating idp: %s", id.Metadata.Name)
		uri := fmt.Sprintf("/auth/v3/sso/idp/%s", id.Metadata.Name)
		_, err := auth.AuthAndRequest(uri, "PUT", id)
		if err != nil {
			return err
		}
	} else {
		log.GetLogger().Debugf("creating idp: %s", id.Metadata.Name)
		uri := "/auth/v3/sso/idp"
		_, err := auth.AuthAndRequest(uri, "POST", id)
		if err != nil {
			return err
		}
	}
	return nil
}
