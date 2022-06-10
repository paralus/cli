package location

import (
	"encoding/json"
	"fmt"

	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/rerror"
	infrav3 "github.com/paralus/paralus/proto/types/infrapb/v3"
)

func ListAllLocation() ([]*infrav3.Metro, error) {
	cfg := config.GetConfig()
	var locations []*infrav3.Metro
	limit := 10000
	b, count, err := ListLocation(cfg.Partner, limit, 0)
	if err != nil {
		return nil, err
	}
	locations = b
	for count > limit {
		offset := limit
		limit = count
		b, _, err = ListLocation(cfg.Partner, limit, offset)
		if err != nil {
			return locations, err
		}
		locations = append(locations, b...)
	}
	return locations, nil
}

/*
ListLocation is used to fetch the locations in the provided project. It accepts a project id
a limit and an offset as inputs. It returns an error if there was a problem while fetching
the locations. The function will return the list of locations, total locations count, and an error
if applicable.
*/
func ListLocation(partner string, limit, offset int) ([]*infrav3.Metro, int, error) {
	// check to make sure the limit or offset is not negative
	if limit < 0 || offset < 0 {
		return nil, 0, fmt.Errorf("provided limit (%d) or offset (%d) cannot be negative", limit, offset)
	}
	auth := config.GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/infra/v3/partner/%s/location", partner)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, 0, rerror.CrudErr{
			Type: "locations",
			Name: "",
			Op:   "list",
		}
	}
	a := infrav3.LocationList{}
	err = json.Unmarshal([]byte(resp), &a)
	if err != nil {
		return nil, -1, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return a.Items, len(a.Items), nil
}

// GetLocation fetches a single location
func GetLocation(locationName string) (*infrav3.Metro, error) {
	ls, err := ListAllLocation()
	if err != nil {
		return nil, err
	}
	for _, l := range ls {
		if l.Name == locationName {
			return l, nil
		}
	}

	return nil, rerror.ResourceNotFound{
		Type: "location",
		Name: locationName,
	}
}
