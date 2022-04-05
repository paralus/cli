package printer

import (
	"bytes"
	"testing"
)

func TestPrintTableJsonPath(t *testing.T) {
	type args struct {
		jsonStrings []string
		columns     []string
		columnPaths map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Simple Table Test",
			args: args{
				jsonStrings: []string{
					`{"name":{"first":"Janet","last":"Prichard"},"age":47}`,
					`{"name":{"first":"Dale","last":"Murphy"},"age":44}`,
				},
				columns: []string{
					"First Name",
					"Last Name",
				},
				columnPaths: map[string]string{
					"First Name": "name.first",
					"Last Name":  "name.last",
				},
			},
			want: `+------------+-----------+
| FIRST NAME | LAST NAME |
+------------+-----------+
| Janet      | Prichard  |
+------------+-----------+
| Dale       | Murphy    |
+------------+-----------+
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			PrintTableJsonPath(tt.args.jsonStrings, tt.args.columns, tt.args.columnPaths, &buf)
			if buf.String() != tt.want {
				t.Error(buf.String())
				t.Fail()
			}
		})
	}
}

func TestPrintTable(t *testing.T) {
	type args struct {
		columns []string
		rows    [][]string
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
	}{
		{
			name: "Simple Table Test",
			args: args{
				columns: []string{
					"First Name",
					"Last Name",
				},
				rows: [][]string{
					{"Janet", "Prichard"},
					{"Dale", "Murphy"},
				},
			},
			wantWriter: `+------------+-----------+
| FIRST NAME | LAST NAME |
+------------+-----------+
| Janet      | Prichard  |
+------------+-----------+
| Dale       | Murphy    |
+------------+-----------+
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			PrintTable(tt.args.columns, tt.args.rows, writer)
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("PrintTable() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
