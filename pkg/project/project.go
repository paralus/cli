package project

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/constants"
	"github.com/paralus/cli/pkg/log"
	"github.com/paralus/cli/pkg/output"
	"github.com/paralus/cli/pkg/prefix"
	"github.com/paralus/cli/pkg/rerror"
	"github.com/paralus/cli/pkg/utils"
	commonv3 "github.com/paralus/paralus/proto/types/commonpb/v3"
	systemv3 "github.com/paralus/paralus/proto/types/systempb/v3"

	"github.com/spf13/cobra"
)

func ListProjectsWithCmd(cmd *cobra.Command) (*systemv3.ProjectList, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/projects", cfg.Partner, cfg.Organization)
	uri = utils.AddPagenationToRequestWithCmd(cmd, uri)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, rerror.CrudErr{
			Type: "project",
			Name: "",
			Op:   "list",
		}
	}
	projs := &systemv3.ProjectList{}
	err = json.Unmarshal([]byte(resp), projs)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return projs, nil
}

func GetProjectByName(projectName string) (*systemv3.Project, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/project/%s", cfg.Partner, cfg.Organization, projectName)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	proj := &systemv3.Project{}
	err = json.Unmarshal([]byte(resp), proj)
	if err != nil {
		return nil, err
	}

	return proj, nil
}

func getCount(obj interface{}) int {
	countInt := 0
	count, err := jsonpath.JsonPathLookup(obj, "$.metadata.count")
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

func newProjectListSpec(obj interface{}) *output.OutputListSpec {
	countInt := getCount(obj)

	log.GetLogger().Debugf("count = %d", countInt)
	columnSpec := []output.ColumnSpec{
		{Header: "NAME", JsonPath: "name"},
	}

	spec := &output.OutputListSpec{
		NRows:   countInt,
		Base:    "$.items[%d]",
		Columns: columnSpec,
	}

	return spec
}

func NewProjectFromResponse(json_data []byte) (*systemv3.Project, error) {
	var pr *systemv3.Project
	if err := json.Unmarshal(json_data, &pr); err != nil {
		return nil, err
	}
	if pr == nil {
		return nil, fmt.Errorf("project could not be found")
	}
	return pr, nil
}

func ListPrint(cmd *cobra.Command, jsonObj []byte) error {
	return output.PrintList(cmd, jsonObj, newProjectListSpec)
}

func Print(cmd *cobra.Command, jsonObj []byte) error {
	p, err := NewProjectFromResponse(jsonObj)
	if err != nil {
		return err
	}

	w := prefix.NewPrefixWriter()
	w.Write(0, "Name: %s\n", p.Metadata.Name)

	w.Flush()

	return nil
}

func GetProjectIdListFromProjectsString(commaSeparatedProjects string) ([]string, error) {
	projectList := strings.Split(commaSeparatedProjects, ",")
	projectIdList := []string{}
	for _, projectName := range projectList {
		pr, err := GetProjectByName(projectName)
		if err != nil {
			return []string{}, fmt.Errorf("project %s was not found", projectName)
		}
		projectIdList = append(projectIdList, pr.Metadata.Id)
	}
	return projectIdList, nil
}

func CreateProject(name, description string) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	err, matched := utils.MatchStringToRegx(name, constants.PROJECT_NAME_REGEX)
	if err != nil || !matched {
		return fmt.Errorf("cluster name contains invalid characters. valid characters are `%s`", constants.PROJECT_NAME_REGEX)
	}

	project := systemv3.Project{
		Kind: "Project",
		Metadata: &commonv3.Metadata{
			Name:        name,
			Description: description,
		},
		Spec: &systemv3.ProjectSpec{
			Default: false,
		},
	}

	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/project", cfg.Partner, cfg.Organization)
	_, err = auth.AuthAndRequest(uri, "POST", project)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProject(project string) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/project/%s", cfg.Partner, cfg.Organization, project)
	_, err := auth.AuthAndRequest(uri, "DELETE", nil)
	return err
}

// Apply project takes the project details and sends it to the core
func ApplyProject(proj *systemv3.Project) error {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()

	projExisting, _ := GetProjectByName(proj.Metadata.Name)
	if projExisting != nil {
		log.GetLogger().Debugf("updating project: %s", proj.Metadata.Name)
		uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/project/%s", cfg.Partner, cfg.Organization, proj.Metadata.Name)
		_, err := auth.AuthAndRequest(uri, "PUT", proj)
		if err != nil {
			return rerror.CrudErr{
				Type: "project",
				Name: proj.Metadata.Name,
				Op:   "update",
			}
		}
	} else {
		log.GetLogger().Debugf("creating project: %s", proj.Metadata.Name)
		uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/project", cfg.Partner, cfg.Organization)
		_, err := auth.AuthAndRequest(uri, "POST", proj)
		if err != nil {
			return rerror.CrudErr{
				Type: "project",
				Name: proj.Metadata.Name,
				Op:   "create",
			}
		}
	}
	return nil
}
