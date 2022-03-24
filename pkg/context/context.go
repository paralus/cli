package context

import (
	"bytes"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/rafaylabs/rcloud-cli/pkg/log"
	"github.com/rafaylabs/rcloud-cli/pkg/utils"
)

type CliContext struct {
	ConfigDir        string `json:"config_dir,omitempty"`
	ConfigFile       string `json:"config_file,omitempty"`
	Verbose          bool   `json:"verbose"`
	Debug            bool   `json:"debug"`
	StructuredOutput bool   `json:"structured_output"`
	V3               bool   `json:"v3"`
}

var context = &CliContext{
	ConfigDir:        filepath.Join(utils.GetUserHome(), ".rafay", "cli"),
	ConfigFile:       "config.json",
	Verbose:          false,
	Debug:            false,
	StructuredOutput: false,
	V3:               false,
}

func GetContext() *CliContext {
	return context
}

func (c *CliContext) WriteReadable(w io.Writer) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")
	enc.Encode(c)
}

func (c *CliContext) Write(w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(c)
}

func (c *CliContext) UseStructuredOutput() bool {
	return c.StructuredOutput
}

func (c *CliContext) ConfigFilename() string {
	return filepath.Join(c.ConfigDir, c.ConfigFile)
}

func (c *CliContext) Log(title string) {
	buf := bytes.NewBufferString(title)
	c.Write(buf)
	log.GetLogger().Infof(strings.TrimSuffix(buf.String(), "\n"))
}

func (c *CliContext) DLog(title string) {
	buf := bytes.NewBufferString(title)
	c.Write(buf)
	log.GetLogger().Debugf(strings.TrimSuffix(buf.String(), "\n"))
}
