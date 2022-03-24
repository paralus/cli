package commands

import (
	"testing"

	"github.com/rafaylabs/rcloud-cli/test"
	"github.com/stretchr/testify/assert"
)

func TestDownloadKubeconfigOptions_ValidateArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "test download kubeconfig",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "test download kubeconfig with arguments",
			args:    []string{"c1", "c2", "c3"},
			wantErr: false,
		},
		{
			name:    "test download kubeconfig with cluster flag",
			args:    []string{"--" + DownloadKubeconfigClusterFlag, "c1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewDownloadKubeconfigOptions(test.NewNoopLogger())
			assert.NotNil(t, o, "DownloadKubeconfigOptions is nil")
			c := newTestWithConfigCmd("c", o.Validate, test.EmptyRun)
			o.AddFlags(c)
			if _, err := test.ExecuteCommand(c, tt.args...); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
