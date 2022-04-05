package config

import (
	"os"
	"testing"

	"github.com/RafayLabs/rcloud-cli/pkg/context"
	"github.com/stretchr/testify/assert"
)

func unsetEnvs() {
	os.Unsetenv("RCTL_PROFILE")
	os.Unsetenv("RCTL_REST_ENDPOINT")
	os.Unsetenv("RCTL_OPS_ENDPOINT")
	os.Unsetenv("RCTL_API_KEY")
	os.Unsetenv("RCTL_API_SECRET")
}

func setEnvs() {
	os.Setenv("RCTL_PROFILE", "test_env_val")
	os.Setenv("RCTL_REST_ENDPOINT", "test_env_val")
	os.Setenv("RCTL_OPS_ENDPOINT", "test_env_val")
	os.Setenv("RCTL_API_KEY", "test_env_val")
	os.Setenv("RCTL_API_SECRET", "test_env_val")
}

func TestInitConfig(t *testing.T) {
	type args struct {
		ctx *context.CliContext
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create config with file",
			args: args{&context.CliContext{
				ConfigDir:  "../../test/data",
				ConfigFile: "config.json",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitConfig(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitConfig_FakeConfigFileNoEnvs(t *testing.T) {
	unsetEnvs()
	fakeConfig := &context.CliContext{
		ConfigDir:  "../../test/data",
		ConfigFile: "config.json",
	}

	if err := InitConfig(fakeConfig); err != nil {
		t.Errorf("InitConfig() error = %v, should be nil", err)
	}

	config := GetConfig()
	configExpected := &Config{
		Profile:             "dev",
		RESTEndpoint:        "localhost:11000",
		OPSEndpoint:         "localhost:11000",
		SkipServerCertValid: "false",
		APIKey:              "7719ae598b735dfd1ebb3c35ab9a13bede0882b0",
		APISecret:           "6430334e3e68bdc815f9559bc16f63695adb120d7e99821f34444127072727ca",
		Partner:             "finman",
		Organization:        "finmanorg",
		Project:             "myproject",
	}

	assert.Equal(t, config, configExpected, "config is not as expected")

	tracker := GetConfigTracker()
	trackerExpected := ConfigTracker{
		Config: *configExpected,
		Source: Config{
			Profile:             "File ../../test/data/config.json",
			RESTEndpoint:        "File ../../test/data/config.json",
			OPSEndpoint:         "File ../../test/data/config.json",
			APIKey:              "File ../../test/data/config.json",
			APISecret:           "File ../../test/data/config.json",
			Project:             "File ../../test/data/config.json",
			Organization:        "File ../../test/data/config.json",
			Partner:             "File ../../test/data/config.json",
			SkipServerCertValid: "File ../../test/data/config.json",
		},
	}

	assert.Equal(t, tracker.Source.Profile, trackerExpected.Source.Profile, "trackerConfig is not as expected")
	assert.Equal(t, tracker.Source.RESTEndpoint, trackerExpected.Source.RESTEndpoint, "trackerConfig is not as expected")
	assert.Equal(t, tracker.Source.OPSEndpoint, trackerExpected.Source.OPSEndpoint, "trackerConfig is not as expected")
	assert.Equal(t, tracker.Source.APIKey, trackerExpected.Source.APIKey, "trackerConfig is not as expected")
	assert.Equal(t, tracker.Source.APISecret, trackerExpected.Source.APISecret, "trackerConfig is not as expected")

}

func TestInitConfig_FakeConfigFileWithEnvs(t *testing.T) {
	setEnvs()
	fakeConfig := &context.CliContext{
		ConfigDir:  "../../test/data",
		ConfigFile: "config.json",
	}

	if err := InitConfig(fakeConfig); err != nil {
		t.Errorf("InitConfig() error = %v, should be nil", err)
	}

	config := GetConfig()
	configExpected := &Config{
		Profile:             "test_env_val",
		RESTEndpoint:        "test_env_val",
		OPSEndpoint:         "test_env_val",
		APIKey:              "test_env_val",
		APISecret:           "test_env_val",
		Partner:             "finman",
		Organization:        "finmanorg",
		Project:             "myproject",
		SkipServerCertValid: "false",
	}

	assert.Equal(t, config, configExpected, "config is not as expected")

	tracker := GetConfigTracker()
	trackerExpected := ConfigTracker{
		Config: *configExpected,
		Source: Config{
			Profile:             "Environment variable",
			RESTEndpoint:        "Environment variable",
			OPSEndpoint:         "Environment variable",
			APIKey:              "Environment variable",
			APISecret:           "Environment variable",
			Project:             "File ../../test/data/config.json",
			Organization:        "File ../../test/data/config.json",
			Partner:             "File ../../test/data/config.json",
			SkipServerCertValid: "File ../../test/data/config.json",
		},
	}

	assert.Equal(t, tracker.Source.Profile, trackerExpected.Source.Profile, "trackerConfig is not as expected")
	assert.Equal(t, tracker.Source.RESTEndpoint, trackerExpected.Source.RESTEndpoint, "trackerConfig is not as expected")
	assert.Equal(t, tracker.Source.OPSEndpoint, trackerExpected.Source.OPSEndpoint, "trackerConfig is not as expected")
	assert.Equal(t, tracker.Source.APIKey, trackerExpected.Source.APIKey, "trackerConfig is not as expected")
	assert.Equal(t, tracker.Source.APISecret, trackerExpected.Source.APISecret, "trackerConfig is not as expected")

}
