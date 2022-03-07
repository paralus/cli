package utils

import (
	"fmt"

	"github.com/spf13/cobra"
)

func addPagenation(offset, limit int, url string) string {
	if limit != 10 || offset != 0 {
		return url + fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	}
	return url
}

// Add pagination flags to a URL request
func AddPagenationToRequestWithCmd(cmd *cobra.Command, url string) string {
	limit, _ := cmd.Flags().GetInt("limit")
	offset, _ := cmd.Flags().GetInt("offset")

	if limit == 0 {
		limit = 10
	}

	return addPagenation(offset, limit, url)
}

// Add pagination flags to a URL request
func AddPagenationToRequest(url string, limit, offset int) (string, error) {
	switch {
	case limit == 0:
		return url, fmt.Errorf("limit cannot be 0")
	case limit < 0:
		return url, fmt.Errorf("limit cannot be less than 0")
	case offset < 0:
		return url, fmt.Errorf("limit cannot be less than 0")
	}

	return addPagenation(offset, limit, url), nil
}

// Full
func FullPagenationToRequest(url string) string {
	return addPagenation(0, 1000000, url)
}

// Get all data from paginated API
func ListAllInstances(
	baseUri string,
	getObjects func(string) ([]interface{}, int, error)) ([]interface{}, int, error) {
	var limit, offset, count int = 50, 0, 0 // default limit is set to 50
	var instanceList []interface{}

	for count >= 0 {
		uri := fmt.Sprintf("%s?&options.limit=%d&options.offset=%d", baseUri, limit, offset)
		// get objects for this iteration
		objects, totalCount, err := getObjects(uri)
		if err != nil {
			return nil, -1, err
		}
		// set count in initial loop if total number of objects for the query > 0
		if totalCount == 0 {
			// no objects found
			return nil, 0, nil
		} else if count == 0 && offset == 0 {
			// set count in initial loop
			count = totalCount
		}
		// append objects from this iteration
		instanceList = append(instanceList, objects...)
		// evaluate remaining count, break if all objects retrieved
		if count > limit {
			count -= limit
			offset += limit
		} else {
			break
		}
	}
	return instanceList, len(instanceList), nil
}
