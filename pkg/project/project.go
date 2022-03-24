package project

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oliveagle/jsonpath"
	"github.com/rafaylabs/rcloud-cli/pkg/config"
	"github.com/rafaylabs/rcloud-cli/pkg/log"
	"github.com/rafaylabs/rcloud-cli/pkg/models"
	"github.com/rafaylabs/rcloud-cli/pkg/output"
	"github.com/rafaylabs/rcloud-cli/pkg/prefix"
	"github.com/rafaylabs/rcloud-cli/pkg/rerror"
	"github.com/rafaylabs/rcloud-cli/pkg/utils"

	"github.com/spf13/cobra"
)

func ListProjectsWithCmd(cmd *cobra.Command) (*models.ProjectList, error) {
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
	projs := &models.ProjectList{}
	err = json.Unmarshal([]byte(resp), projs)
	if err != nil {
		return nil, fmt.Errorf("there was an error while unmarshalling: %v", err)
	}
	return projs, nil
}

func GetProjectByName(projectName string) (*models.Project, error) {
	cfg := config.GetConfig()
	auth := cfg.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/project/%s", cfg.Partner, cfg.Organization, projectName)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return nil, err
	}
	proj := &models.Project{}
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

func NewProjectFromResponse(json_data []byte) (*models.Project, error) {
	var pr *models.Project
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

	project := models.Project{
		Kind: "Project",
		Metadata: models.Metadata{
			Name:        name,
			Description: description,
		},
		Spec: models.ProjectSpec{
			Default: false,
		},
	}

	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/project", cfg.Partner, cfg.Organization)
	_, err := auth.AuthAndRequest(uri, "POST", project)
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
func ApplyProject(proj *models.Project) error {
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
