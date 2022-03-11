package utils

import (
	"testing"
)

func TestFileExists(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test file exists",
			args: args{
				filename: "../testdata/addon/yaml/elastic-op.yml",
			},
			want: true,
		},
		{
			name: "Test file not exists",
			args: args{
				filename: "../testdata/addon/yaml/1234.yml",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileExists(tt.args.filename); got != tt.want {
				t.Errorf("FileExists() = %v, want %v", got, tt.want)
			}
		})
	}
}
