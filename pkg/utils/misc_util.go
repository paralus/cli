package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

type ErrorDetails struct {
	ErrorCode string `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	Detail    string `json:"detail,omitempty" yaml:"detail,omitempty"`
	Info      string `json:"info,omitempty" yaml:"info,omitempty"`
}

type RafayErrorMessage struct {
	StatusCode int            `json:"status_code,omitempty" yaml:"status_code,omitempty"`
	Details    []ErrorDetails `json:"details,omitempty" yaml:"details,omitempty"`
}

func GetUserHome() string {
	homeEnvVariable := "HOME"
	if runtime.GOOS == "windows" {
		homeEnvVariable = "USERPROFILE"
	}
	return os.Getenv(homeEnvVariable)
}

func FormatYamlMessage(data interface{}) (string, error) {
	var ret string
	bArr, err := yaml.Marshal(data)
	if err == nil {
		ret = string(bArr)
	}
	return ret, err
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func FullPath(parentFile, path string) string {
	if path == "" || filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(filepath.Dir(parentFile), path)
}

func FullPaths(parentFile, path string) string {
	allPaths := strings.Split(path, ",")
	if len(allPaths) <= 1 {
		return FullPath(parentFile, path)
	}
	allFullPaths := make([]string, len(allPaths))
	for i, aPath := range allPaths {
		allFullPaths[i] = FullPath(parentFile, aPath)
	}
	return strings.Join(allFullPaths, ",")
}

func GetAsString(i interface{}) string {
	if i == nil {
		return ""
	}
	return i.(string)
}

func GetAsMap(array []string) map[string]string {
	asMap := make(map[string]string)
	for _, entry := range array {
		asMap[entry] = ""
	}
	return asMap
}

func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func PrettyPrint(responseStr string) {
	if !gjson.Valid(responseStr) {
		fmt.Println(responseStr)
		return
	}
	result := gjson.Get(responseStr, "@pretty")
	fmt.Println(result.String())
}

func ExpandFile(path string) (string, error) {
	if len(path) == 0 || path[0] != '~' {
		abs, err := filepath.Abs(path)
		return abs, err
	}

	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, path[1:]), nil
}

func ProcessValueFiles(valuesFilePath string) (string, error) {
	allValuesFileNames := ""
	allValuesFiles := strings.Split(valuesFilePath, ",")
	for _, valuesFile := range allValuesFiles {
		if valuesFile != "" {
			absFile, err := ExpandFile(valuesFile)
			if err != nil {
				return "", fmt.Errorf("values file %s does not exist error %s", valuesFile, err.Error())
			}
			if !FileExists(absFile) {
				return "", fmt.Errorf("values file %s does not exist", valuesFile)
			}
			allValuesFileNames = allValuesFileNames + absFile + ","
		}
	}
	return allValuesFileNames, nil
}
