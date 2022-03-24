package config

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/rafaylabs/rcloud-cli/pkg/authprofile"

	"github.com/rafaylabs/rcloud-cli/pkg/context"

	oruntime "github.com/go-openapi/runtime"
	oclient "github.com/go-openapi/runtime/client"

	"github.com/rafaylabs/rcloud-cli/pkg/log"
)

var logger = log.GetLogger()

type Config struct {
	Profile             string `json:"profile,omitempty"`
	SkipServerCertValid string `json:"skip_server_cert_check,omitempty"`
	RESTEndpoint        string `json:"rest_endpoint,omitempty"`
	OPSEndpoint         string `json:"ops_endpoint,omitempty"`
	APIKey              string `json:"api_key,omitempty"`
	APISecret           string `json:"api_secret,omitempty"`
	Partner             string `json:"partner,omitempty"`
	Organization        string `json:"organization,omitempty"`
	Project             string `json:"project,omitempty"`
}

type ConfigTracker struct {
	Config
	Source Config `json:"source"`
}

var config ConfigTracker = ConfigTracker{}

func NewProductionConfig() *Config {
	return &Config{
		Profile:             "prod",
		RESTEndpoint:        "console.rafay.dev",
		OPSEndpoint:         "ops.rafay.dev",
		SkipServerCertValid: "false",
	}
}

func NewStageConfig() *Config {
	return &Config{
		Profile:             "staging",
		RESTEndpoint:        "console.stage.rafay.dev",
		OPSEndpoint:         "ops.stage.rafay.dev",
		SkipServerCertValid: "false",
	}
}

func NewOpDevConfig() *Config {
	return &Config{
		Profile:             "opdev",
		RESTEndpoint:        "localhost:80",
		OPSEndpoint:         "localhost:80",
		SkipServerCertValid: "false",
	}
}

func NewDefaultConfig(profile string) *Config {
	switch profile {
	case "dev":
		return &Config{}
	case "staging":
		return NewStageConfig()
	case "opdev":
		return NewOpDevConfig()
	case "prod":
		fallthrough
	case "":
		return NewProductionConfig()
	}

	return nil

}

// Return the merged context *b into the *c. When a field in *c is
// not defined, it is merged in from *b.
func (c *ConfigTracker) Merge(b *Config, source string, override bool) {
	// merge if empty values, or if override flag is set
	if len(c.Profile) == 0 || (override && len(b.Profile) != 0) {
		c.Profile = b.Profile
		c.Source.Profile = source
	}

	if len(c.RESTEndpoint) == 0 || (override && len(b.RESTEndpoint) != 0) {
		c.RESTEndpoint = b.RESTEndpoint
		c.Source.RESTEndpoint = source
	}

	if len(c.OPSEndpoint) == 0 || (override && len(b.OPSEndpoint) != 0) {
		c.OPSEndpoint = b.OPSEndpoint
		c.Source.OPSEndpoint = source
	}

	if len(c.APIKey) == 0 || (override && len(b.APIKey) != 0) {
		c.APIKey = b.APIKey
		c.Source.APIKey = source
	}

	if len(c.APISecret) == 0 || (override && len(b.APISecret) != 0) {
		c.APISecret = b.APISecret
		c.Source.APISecret = source
	}

	if len(c.Project) == 0 || (override && len(b.Project) != 0) {
		c.Project = b.Project
	}

	if len(c.Organization) == 0 || (override && len(b.Organization) != 0) {
		c.Organization = b.Organization
	}

	if len(c.Partner) == 0 || (override && len(b.Partner) != 0) {
		c.Partner = b.Partner
	}

	if len(c.SkipServerCertValid) == 0 || (override && len(b.SkipServerCertValid) != 0) {
		c.SkipServerCertValid = b.SkipServerCertValid
		c.Source.SkipServerCertValid = source
	}

}

func (c *Config) MiniCheck() error {
	if len(c.Profile) == 0 {
		return fmt.Errorf("profile name not defined")
	}

	// Check if profile is known to us
	if p := NewDefaultConfig(c.Profile); p == nil {
		return fmt.Errorf("unknown profile name")
	}

	if len(c.APIKey) == 0 {
		return fmt.Errorf("api key not defined")
	}

	if len(c.APISecret) == 0 {
		return fmt.Errorf("api secret not defined")
	}

	if len(c.Partner) == 0 {
		return fmt.Errorf("partner not defined")
	}

	return nil
}

func InitConfig(ctx *context.CliContext) error {
	path := ctx.ConfigFilename()
	tmp := &Config{}
	log.GetLogger().Debugf("Config path: %v", path)
	err := tmp.Load(path)
	configTracker := GetConfigTracker()
	if err == nil {
		configTracker.Merge(tmp, fmt.Sprintf("File %s", path), false)
	}
	// configTracker.SkipServerCertValid = "false"

	fromEnv := &Config{
		Profile:             os.Getenv("RCTL_PROFILE"),
		RESTEndpoint:        os.Getenv("RCTL_REST_ENDPOINT"),
		OPSEndpoint:         os.Getenv("RCTL_OPS_ENDPOINT"),
		APIKey:              os.Getenv("RCTL_API_KEY"),
		APISecret:           os.Getenv("RCTL_API_SECRET"),
		SkipServerCertValid: os.Getenv("RCTL_SKIP_SERVER_CERT_VALIDATION"),
	}
	configTracker.Merge(fromEnv, "Environment variable", true)

	// Figure out which profile to use
	profile := configTracker.Profile
	source := ""
	if len(profile) == 0 {
		profile = "prod"
		source = "Default"
	}

	if len(source) != 0 {
		configTracker.Profile = profile
		configTracker.Source.Profile = source
	}

	if p := NewDefaultConfig(profile); p != nil {
		configTracker.Merge(p, fmt.Sprintf("Profile %s default", profile), false)
	}

	return nil
}

func GetConfig() *Config {
	return &config.Config
}

func GetConfigTracker() *ConfigTracker {
	return &config
}

func (c *Config) WriteReadable(w io.Writer) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")
	enc.Encode(c)
}

func (c *Config) Write(w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(c)
}

func (c *Config) Save(filename string) error {
	parentDir := filepath.Dir(filename)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		dirCreateError := os.MkdirAll(parentDir, os.ModePerm)
		if dirCreateError != nil {
			return fmt.Errorf("RCTL config director path [%s] doesn't exist and attempt to create it failed with error [%s]", parentDir, err.Error())
		}
	}
	f, err := os.Create(filename)
	if err != nil {
		log.GetLogger().Debugf("Failed to create file %s\n", filename)
		return err
	}

	w := bufio.NewWriter(f)
	c.WriteReadable(w)
	w.Flush()
	f.Sync()
	f.Close()
	return nil
}

func (c *Config) Load(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(file), &c)
}

func (c *Config) Log(title string) {
	buf := bytes.NewBufferString(title)
	/* Redact API key info in logs */
	logCfg := *c
	logCfg.APIKey = ""
	logCfg.APISecret = ""
	logCfg.Write(buf)
	log.GetLogger().Infof(strings.TrimSuffix(buf.String(), "\n"))
}

func (c *Config) DLog(title string) {
	buf := bytes.NewBufferString(title)
	/* Redact API key info in logs */
	logCfg := *c
	logCfg.APIKey = ""
	logCfg.APISecret = ""
	logCfg.Write(buf)
	log.GetLogger().Debugf(strings.TrimSuffix(buf.String(), "\n"))
}

func (c *Config) OpenAPIAuthInfo() oruntime.ClientAuthInfoWriter {
	return oclient.APIKeyAuth("X-RAFAY-API-KEYID", "header", c.APIKey)
}

func endpointURL(endpoint string) string {
	host, port, err := net.SplitHostPort(endpoint)
	if err == nil && port == "80" {
		return fmt.Sprintf("http://%s", host)
	}

	if err == nil && port == "11000" {
		return fmt.Sprintf("http://%s:%s", host, port)
	}

	return fmt.Sprintf("https://%s", endpoint)
}

func (c *Config) GetAppAuthProfile() *authprofile.Profile {
	var skip bool
	if c.SkipServerCertValid == "true" {
		skip = true
	}
	return &authprofile.Profile{
		Name:                "app",
		Key:                 c.APIKey,
		Secret:              c.APISecret,
		URL:                 endpointURL(c.RESTEndpoint),
		SkipServerCertValid: skip,
	}
}

func (c *Config) GetOpsAuthProfile() *authprofile.Profile {
	profile := c.GetAppAuthProfile()

	profile.Name = "ops"
	profile.URL = endpointURL(c.OPSEndpoint)

	return profile
}

func (c *Config) Output() {
	projectName, err := GetProjectNameById(c.Partner, c.Organization, c.Project)
	if err != nil {
		fmt.Println("Failed to verify project referred by the config. Please make sure API is valid and/or use \"rctl config set project <project name>\" to set valid Project.")
		projectName = "[Error: invalid]"
	}
	fmt.Printf("%-39s %64s\n", "Profile:", c.Profile)
	fmt.Printf("%-39s %64s\n", "REST Endpoint:", c.RESTEndpoint)
	fmt.Printf("%-39s %64s\n", "OPS Endpoint:", c.OPSEndpoint)
	fmt.Printf("%-39s %64s\n", "API Key:", c.APIKey)
	fmt.Printf("%-39s %64s\n", "API Secret:", c.APISecret)
	fmt.Printf("%-39s %64s\n", "Skip Server Certificate Validation:", c.SkipServerCertValid)
	fmt.Printf("%-39s %64s\n", "Project:", projectName)
}

func (t *ConfigTracker) Output() {
	projectName, err := GetProjectNameById(t.Partner, t.Organization, t.Project)
	if err != nil {
		fmt.Println("Project referred by the config doesn't exist or you don't have access to it. Please use \"rctl config set project <project name>\" to set valid Project.")
		projectName = "[Error: invalid]"
	}

	fmt.Printf(
		`Profile:
  From:    %s
  Value:   %s

REST Endpoint:
  From:    %s
  Value:   %s

OPS Endpoint:
  From:    %s
  Value:   %s

API Key:
  From:    %s
  Value:   %s

API Secret:
  From:    %s
  Value:   %s

Skip Server Certificate Validation:
  From:    %s
  Value:   %s

Project:
  From:    %s
  Value:   %s
`,
		t.Source.Profile,
		t.Profile,
		t.Source.RESTEndpoint,
		t.RESTEndpoint,
		t.Source.OPSEndpoint,
		t.OPSEndpoint,
		t.Source.APIKey,
		t.APIKey,
		t.Source.APISecret,
		t.APISecret,
		t.Source.SkipServerCertValid,
		t.SkipServerCertValid,
		t.Source.Project,
		projectName,
	)
}
