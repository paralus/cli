package config

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/RafaySystems/rcloud-cli/pkg/constants"
	"github.com/RafaySystems/rcloud-cli/pkg/exit"
	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"github.com/oliveagle/jsonpath"
	"github.com/segmentio/encoding/json"
	"github.com/spf13/cobra"
)

func refCheck(ids interface{}) (string, error) {
	if idsArray, ok := ids.([]interface{}); ok {
		if len(idsArray) == 1 {
			if id, ok := idsArray[0].(string); ok {
				return id, nil
			} else {
				return "", fmt.Errorf("ref id is not of a string type")
			}
		} else {
			return "", fmt.Errorf("ref object matches more than one element")
		}
	}
	return "", fmt.Errorf("ref object is not an array")
}

func GetProjectIdFromFlagAndConfig(cmd *cobra.Command) (string, error) {
	projectName, err := cmd.Flags().GetString("project")
	if err != nil {
		exit.SetExitWithError(err, "Invalid flag/argument passed for project")
		return "", err
	}
	project := projectName
	if projectName != "" {
		_, err = GetProjectIdByName(projectName)
		if err != nil {
			exit.SetExitWithError(err, "Project not found")
			return "", err
		}
	}
	if project == "" {
		project = GetConfig().Project
	}
	if project == "" {
		return "", fmt.Errorf("project context couldn't be determined. Please use --project argument or init rctl with the project context using \"rctl config set project <project name>\"")
	}
	return project, nil
}

func GetProjectIdFromFlag(cmd *cobra.Command) (string, error) {
	projectName, err := cmd.Flags().GetString("project")
	if err != nil {
		exit.SetExitWithError(err, "Invalid flag/argument passed for project")
		return "", err
	}
	var projectId string
	if projectName != "" {
		projectId, err = GetProjectIdByName(projectName)
		if err != nil {
			exit.SetExitWithError(err, "Project not found")
			return "", err
		}
	}
	if projectId == "" {
		return "", fmt.Errorf("project context couldn't be determined. Please use --project argument or init rctl with the project context using \"rctl config set project <project name>\"")
	}
	return projectId, nil
}

func GetProjectIdByName(name string) (string, error) {
	return GetProjectIdByNameInConfig(GetConfig(), name)
}

func GetProjectIdByNameInConfig(config *Config, name string) (string, error) {
	params := url.Values{}
	params.Add("q", name)
	auth := config.GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner/%s/organization/%s/project/%s", config.Partner, config.Organization, name)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return "", err
	}
	var obj interface{}
	err = json.Unmarshal([]byte(resp), &obj)
	if err != nil {
		return "", err
	}
	filter := fmt.Sprintf("$.results[?(@.name == '%s')].id", name)
	res, err := jsonpath.JsonPathLookup(obj, filter)
	if err == nil {
		id, err := refCheck(res)
		if err != nil {
			log.GetLogger().Debugf("ref check failed %s", err)
		} else {
			return id, nil
		}
	} else {
		log.GetLogger().Debugf("ref looked up name failed %s", err)
	}
	return "", fmt.Errorf("failed to find project with %s", name)
}

func GetProjectNameById(partner, organization, id string) (string, error) {
	auth := GetConfig().GetAppAuthProfile()
	uri := fmt.Sprintf("/auth/v3/partner%s/organization/%s/project/%s/", partner, organization, id)
	resp, err := auth.AuthAndRequest(uri, "GET", nil)
	if err != nil {
		return "", fmt.Errorf("can not find project with id: %s", id)
	}
	var unmarshalledMap map[string]interface{}
	err = json.Unmarshal([]byte(resp), &unmarshalledMap)
	if err != nil {
		return "", err
	}
	if name, exists := unmarshalledMap["name"]; exists {
		return name.(string), nil
	}
	return "", fmt.Errorf("can not find the project name for id: %s", id)
}

func GetProjectNameFromFlagAndConfig(cmd *cobra.Command) (string, error) {
	projectName, err := cmd.Flags().GetString("project")
	if err != nil {
		exit.SetExitWithError(err, "Invalid flag/argument passed for project")
		return "", err
	}
	return projectName, nil
}

func GetV3Context(cmd *cobra.Command) context.Context {
	var isDebug bool
	var err error
	ctx := context.Background()
	if cmd != nil {
		isDebug, err = cmd.Flags().GetBool(constants.DEBUG_FLAG_NAME)
		if err != nil {
			isDebug = false
		}
	} else {
		tflog := os.Getenv("TF_LOG")
		if tflog == "TRACE" || tflog == "DEBUG" {
			isDebug = true
		}
	}

	if isDebug {
		ctx = context.WithValue(ctx, "debug", "true")
	} else {
		ctx = context.WithValue(ctx, "debug", "false")
	}
	return ctx
}
