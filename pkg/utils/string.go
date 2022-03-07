package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"gopkg.in/yaml.v3"
)

func StringOrNone(name string) string {
	if len(name) == 0 {
		return "<none>"
	}

	return name
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func StringInSliceCaseInsensitive(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}

func StringsAreEqualCaseInsensitive(a, b string) bool {
	return strings.ToLower(b) == strings.ToLower(a)
}

// StringSliceContainsDuplicates checks if the string slice contains duplicates,
// and returns the first duplicate, if there are multiple
func StringSliceContainsDuplicates(strSlice []string) (string, bool) {
	keys := make(map[string]bool)
	for _, entry := range strSlice {
		if _, ok := keys[entry]; !ok {
			keys[entry] = true
		} else {
			return entry, true
		}
	}
	return "", false
}

type yamlResourceType struct {
	Kind string `yaml:"kind"`
}

// splitYAML function splits an YAML into multiple YAMLs
func SplitYAML(src []byte) ([][]byte, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(src))
	var dst [][]byte
	for {
		// Read the value into an empty interface
		var temp interface{}
		if err := decoder.Decode(&temp); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		// Convert the interface back to YAML
		yamlPartBytes, err := yaml.Marshal(temp)
		if err != nil {
			return nil, err
		}
		dst = append(dst, yamlPartBytes)
	}
	return dst, nil
}

func JoinYAML(srcs [][]byte) ([]byte, error) {
	var dst, prev []byte
	for i, src := range srcs {
		// Validate YAML
		var temp interface{}
		if err := yaml.Unmarshal(src, &temp); err != nil {
			return nil, fmt.Errorf("Error decoding YAML content")
		}
		if i > 0 {
			sep := "---"
			// If a newline does not exist at the end of previous YAML add one.
			if prev[len(prev)-1] != byte('\n') {
				sep = "\n" + sep
			}
			// If a newline does not exist at the start of current YAML add one.
			if src[0] != byte('\n') {
				sep = sep + "\n"
			}
			dst = append(dst, []byte(sep)...)
		}
		dst = append(dst, src...)
		prev = src
	}
	return dst, nil
}

// SplitYaml splits multiple yaml contents by given delimiter and returns map of kind to list of yaml(string).
func SplitYamlAndGetListByKind(src []byte) (map[string][][]byte, error) {
	// Split YAML into parts
	yamlParts, err := SplitYAML(src)
	if err != nil {
		return nil, err
	}
	// Enumerate and return by kind
	m := make(map[string][][]byte)
	for _, resource := range yamlParts {
		var y yamlResourceType
		err := yaml.Unmarshal(resource, &y)
		if err != nil {
			return nil, err
		}
		if len(y.Kind) > 0 {
			m[y.Kind] = append(m[y.Kind], resource)
		}
	}
	log.GetLogger().Debug("YAML files grouped by kind")
	for kind, list := range m {
		log.GetLogger().Debugf(kind)
		for index, config := range list {
			log.GetLogger().Debugf("\t%d: %#v\n", index, string(config))
		}
	}
	return m, nil
}

func ReadYAMLFileContents(filePath string) ([]byte, error) {
	// check if the file exists
	if !FileExists(filePath) {
		return nil, fmt.Errorf("file %s does not exist", filePath)
	}

	// make sure the file is a yaml file
	if filepath.Ext(filePath) != ".yml" && filepath.Ext(filePath) != ".yaml" {
		return nil, fmt.Errorf("file must a yaml file, file type is %s", filepath.Ext(filePath))
	}

	// open the file
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the file
	fileBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	log.GetLogger().Debugf("YAML File Contents: %s: %#v", filePath, string(fileBytes))
	return fileBytes, nil
}
