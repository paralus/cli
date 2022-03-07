package config

import (
	"testing"

	"github.com/RafaySystems/rcloud-cli/pkg/context"
)

func TestGetProjectIdByName(t *testing.T) {
	// set up the config
	t.Log("Setting up config")
	cliCtx := context.GetContext()
	err := InitConfig(cliCtx)
	if err != nil {
		t.Errorf("Error setting up config: %v", err)
	}
	_, err = GetProjectIdByName("defaultproject")
	if err != nil {
		t.Errorf("GetProjectIdByName() error = %v, should be nil", err)
		return
	}
}
