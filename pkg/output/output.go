package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/paralus/cli/pkg/log"

	"github.com/ghodss/yaml"
	"github.com/liggitt/tabwriter"
	"github.com/oliveagle/jsonpath"
	"github.com/spf13/cobra"
)

type ColumnSpec struct {
	Header   string
	JsonPath string
}

type OutputListSpec struct {
	NRows   int
	Base    string
	Columns []ColumnSpec
}

type MatchSpec struct {
	JsonPath string
	Value    string // Value should match against the value stored in JsonPath
}

type OutputObjectSpec struct {
	NRows int
	Base  string
	Match MatchSpec
}

type Outputer interface {
	Output()
}

func isJson(cmd *cobra.Command) bool {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return false
	}

	return output == "json" || output == "j"
}

func isYaml(cmd *cobra.Command) bool {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return false
	}

	return output == "yaml" || output == "y"
}

func PrintList(cmd *cobra.Command, jsonObj []byte, fp func(obj interface{}) *OutputListSpec) error {
	var obj interface{}
	err := json.Unmarshal(jsonObj, &obj)
	if err != nil {
		return err
	}

	if isJson(cmd) || fp == nil {
		prettyJson, err := json.MarshalIndent(obj, "", "   ")
		if err != nil {
			return err
		}

		fmt.Println(string(prettyJson))
		return nil
	} else if isYaml(cmd) {
		PrintYaml(obj)
		return nil
	}

	spec := fp(obj)
	log.GetLogger().Debugf("spec: %v", spec)

	PrintRows(spec, obj)
	return nil
}

func PrintSelectObject(cmd *cobra.Command, jsonObj []byte, sel string, fp func(obj interface{}, sel string) *OutputObjectSpec) error {
	var obj interface{}
	err := json.Unmarshal(jsonObj, &obj)
	if err != nil {
		return err
	}

	spec := fp(obj, sel)
	log.GetLogger().Debugf("spec: %v", spec)

	for r := 0; r < spec.NRows; r++ {
		base := fmt.Sprintf(spec.Base, r)
		path := fmt.Sprintf("%s.%s", base, spec.Match.JsonPath)
		res, err := jsonpath.JsonPathLookup(obj, path)
		if err != nil {
			continue
		}

		if value, ok := res.(string); ok != true || value != spec.Match.Value {
			continue
		}

		// Matching object found. Print it.
		obj, _ := jsonpath.JsonPathLookup(obj, base)
		PrintJson(obj, true)
		break
	}

	return nil
}

func PrintObject(cmd *cobra.Command, jsonObj []byte, fp func(obj interface{}) string) error {
	var obj interface{}
	err := json.Unmarshal(jsonObj, &obj)
	if err != nil {
		return err
	}

	if isJson(cmd) || fp == nil {
		PrintJson(obj, true)
		return nil
	} else if isYaml(cmd) {
		PrintYaml(obj)
		return nil
	}

	fmt.Println(fp(obj))

	return nil
}

func PrintRow(spec *OutputListSpec, obj interface{}, row int, w *tabwriter.Writer) {
	// if row == -1, print the header
	if row == -1 {
		for c := 0; c < len(spec.Columns); c++ {
			fmt.Fprintf(w, "%s\t", spec.Columns[c].Header)
		}
	} else {
		base := fmt.Sprintf(spec.Base, row)
		for c := 0; c < len(spec.Columns); c++ {
			path := fmt.Sprintf("%s.%s", base, spec.Columns[c].JsonPath)
			res, err := jsonpath.JsonPathLookup(obj, path)
			if err != nil {
				fmt.Fprintf(w, " \t")
			} else {
				fmt.Fprintf(w, "%v\t", res)
			}
		}
	}
	fmt.Fprintln(w)
}

func PrintRows(spec *OutputListSpec, obj interface{}) {
	n := spec.NRows
	if n <= 0 {
		return
	}

	w := new(tabwriter.Writer)
	// Format left aligned in space-separated colums of minimal width 5
	// and at least one blank of padding (so wider column entries do not
	// touch each other).
	w.Init(os.Stdout, 3, 0, 3, ' ', tabwriter.RememberWidths)

	// Print the headers
	PrintRow(spec, obj, -1, w)

	for r := 0; r < n; r++ {
		PrintRow(spec, obj, r, w)
	}

	w.Flush()
}

func PrintJson(obj interface{}, pretty bool) {
	enc := json.NewEncoder(os.Stdout)
	if pretty == true {
		enc.SetIndent("", "  ")
	}
	enc.Encode(obj)
}

func PrintYaml(obj interface{}) {
	d, err := yaml.Marshal(obj)
	if err != nil {
		return
	}

	fmt.Printf("%s", string(d))
}

func PrintOutputer(cmd *cobra.Command, obj Outputer) {
	if isJson(cmd) {
		PrintJson(obj, true)
	} else if isYaml(cmd) {
		PrintYaml(obj)
	} else {
		obj.Output()
	}
}

func PrintJsonObject(jsonObj []byte) error {
	var obj interface{}
	err := json.Unmarshal(jsonObj, &obj)
	if err != nil {
		return err
	}

	PrintJson(obj, true)
	return nil
}

func PrintYamlObject(jsonObj []byte) error {
	yaml, err := yaml.JSONToYAML(jsonObj)
	if err != nil {
		return err

	}
	fmt.Println(string(yaml))
	return nil
}
