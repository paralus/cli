package printer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v2"
)

func newTable(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	return table
}

/*
PrintTableJsonPath will print a table with the supplied JSON data and map into the supplied Writer.
The map contains the header and the json path to use access the value.
The table printed will be something like:
+---------------+---------------+
| COLUMN NAME 1 | COLUMN NAME 2 |
+---------------+---------------+
| Some data 1   | Some data 1   |
| Some data 2   | Some data 2   |
+-------------------+-----------+
The function accepts a gjson query for columnsPaths instead of data, so it can query the each json string passed in
*/
func PrintTableJsonPath(jsonStrings []string, columns []string, columnsPaths map[string]string, writer io.Writer) {
	table := newTable(writer)
	// set the headers
	table.SetHeader(columns)

	// add data
	for _, v := range jsonStrings {
		rowData := make([]string, 0, len(columns))
		for _, c := range columns {
			rowData = append(rowData, gjson.Get(v, columnsPaths[c]).String())
		}
		table.Append(rowData)
	}

	table.Render() // Send output
}

func PrintTable(columns []string, rows [][]string, writer io.Writer) {
	table := newTable(writer)
	table.SetRowLine(true)
	// set the headers
	table.SetHeader(columns)
	// add data
	for _, row := range rows {
		table.Append(row)
	}

	table.Render() // Send output
}

// PrintOutputJsonOrYaml will print json or yaml to writer based on output value
func PrintOutputJsonOrYaml(output string, i interface{}, writer io.Writer) (bool, error) {
	var err error
	switch output {
	case "json":
		err = PrintJson(i, writer)
	case "yaml":
		err = PrintYaml(i, writer)
	default:
		return false, fmt.Errorf("output value must be json or yaml")
	}
	return true, err
}

func PrintJson(i interface{}, writer io.Writer) error {
	en := json.NewEncoder(writer)
	en.SetIndent("", " ")
	return en.Encode(i)
}

func PrintYaml(i interface{}, writer io.Writer) error {
	en := yaml.NewEncoder(writer)
	return en.Encode(i)
}
