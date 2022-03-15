package oidc

import (
	"encoding/json"
	"fmt"

	"github.com/RafaySystems/rcloud-cli/pkg/config"
	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"github.com/RafaySystems/rcloud-cli/pkg/models"
	"github.com/RafaySystems/rcloud-cli/pkg/output"
	"github.com/RafaySystems/rcloud-cli/pkg/prefix"
	"github.com/RafaySystems/rcloud-cli/pkg/rerror"
	"github.com/RafaySystems/rcloud-cli/pkg/utils"
	"github.com/oliveagle/jsonpath"
	"github.com/spf13/cobra"
)

func ListOIDCWithCmd(cmd *cobra.Command) (*models.OIDCProviderList, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := "/auth/v3/sso/oidc/provider"
	uri = utils.AddPagenationToRequestWithCmd(cmd, uri)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, rerror.CrudErr{
			Type: "user",
			Name: "",
			Op:   "list",
		}
	}
	oidcs := &models.OIDCProviderList{}
	err = json.Unmarshal([]byte(resp), oidcs)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return oidcs, nil
}

func GetOIDCProviderByName(name string) (*models.OIDCProvider, error) {
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/sso/oidc/provider/%s", name)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	oidcProvider := &models.OIDCProvider{}
	err = json.Unmarshal([]byte(resp), oidcProvider)
	if err != nil {
		return nil, err
	}

	return oidcProvider, nil

}

func NewOIDCFromResponse(json_data []byte) (*models.OIDCProvider, error) {
	var ur models.OIDCProvider
	if err := json.Unmarshal(json_data, &ur); err != nil {
		return nil, err
	}
	fmt.Println(ur)
	return &ur, nil
}

func CreateOIDCProvider(oidc *models.OIDCProvider) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	//set partner and organization
	oidc.Metadata.Partner = cfg.Partner
	oidc.Metadata.Organization = cfg.Organization

	uri := "/auth/v3/sso/oidc/provider"
	_, err := auth.AuthAndRequest(uri, "POST", oidc)
	if err != nil {
		return err
	}
	return nil
}

func DeleteOIDCProvider(name string) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/sso/oidc/provider/%s", name)
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
	g, err := NewOIDCFromResponse(jsonObj)
	if err != nil {
		return err
	}

	w := prefix.NewPrefixWriter()
	w.Write(0, "Name: %s\n", g.Metadata.Name)
	w.Write(0, "ClientID: %s\n", g.Spec.ClientId)
	w.Write(0, "Provider Name: %s\n", g.Spec.ProviderName)
	w.Write(0, "Mapper URL: %s\n", g.Spec.MapperUrl)
	w.Write(0, "Callback URL: %s\n", g.Spec.CallbackUrl)

	w.Flush()

	return nil
}

// Apply oidc takes the idp details and sends it to the core
func ApplyOIDC(id *models.OIDCProvider) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	idpExisting, _ := GetOIDCProviderByName(id.Metadata.Name)
	if idpExisting != nil {
		log.GetLogger().Debugf("updating idp: %s", id.Metadata.Name)
		//set partner and organization
		id.Metadata.Partner = cfg.Partner
		id.Metadata.Organization = cfg.Organization
		uri := fmt.Sprintf("/auth/v3/sso/oidc/provider/%s", id.Metadata.Name)
		_, err := auth.AuthAndRequest(uri, "PUT", id)
		if err != nil {
			return rerror.CrudErr{
				Type: "idp",
				Name: id.Metadata.Name,
				Op:   "update",
			}
		}
	} else {
		log.GetLogger().Debugf("creating idp: %s", id.Metadata.Name)
		//set partner and organization
		id.Metadata.Partner = cfg.Partner
		id.Metadata.Organization = cfg.Organization
		uri := "/auth/v3/sso/oidc/provider"
		_, err := auth.AuthAndRequest(uri, "POST", id)
		if err != nil {
			return rerror.CrudErr{
				Type: "idp",
				Name: id.Metadata.Name,
				Op:   "create",
			}
		}
	}
	return nil
}
