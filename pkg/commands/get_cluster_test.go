package commands

import (
	"testing"

	"github.com/RafayLabs/rcloud-cli/test"
	"github.com/stretchr/testify/assert"
)

func TestGetClusterOptions_ValidateArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "test list cluster",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "test get cluster",
			args:    []string{"c1"},
			wantErr: false,
		},
		{
			name:    "test get multiple cluster",
			args:    []string{"c1", "c2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewGetClusterOptions(test.NewNoopLogger())
			assert.NotNil(t, o, "GetClusterOptions is nil")
			c := newTestWithConfigCmd("c", o.Validate, test.EmptyRun)
			o.AddFlags(c)
			if _, err := test.ExecuteCommand(c, tt.args...); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
