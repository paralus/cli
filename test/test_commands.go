package test

import (
	"bytes"
	ctx "context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/context"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func EmptyRun(*cobra.Command, []string) error { return nil }

func ExecuteCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, output, err = executeCommandC(root, args...)
	return output, err
}

func ExecuteCommandWithContext(ctx ctx.Context, root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.ExecuteContext(ctx)

	return buf.String(), err
}

func executeCommandC(root *cobra.Command, args ...string) (c *cobra.Command, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	c, err = root.ExecuteC()

	return c, buf.String(), err
}

func ResetCommandLineFlagSet() {
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
}

func CheckStringContains(t *testing.T, got, expected string) {
	if !strings.Contains(got, expected) {
		t.Errorf("Expected to contain: \n %v\nGot:\n %v\n", expected, got)
	}
}

func CheckStringOmits(t *testing.T, got, expected string) {
	if strings.Contains(got, expected) {
		t.Errorf("Expected to not contain: \n %v\nGot: %v", expected, got)
	}
}

// logger used for testing, does nothing
type NoopLogger struct{}

func (n NoopLogger) Debug(args ...interface{}) {}

func (n NoopLogger) Info(args ...interface{}) {}

func (n NoopLogger) Warn(args ...interface{}) {}

func (n NoopLogger) Error(args ...interface{}) {}

func (n NoopLogger) DPanic(args ...interface{}) {}

func (n NoopLogger) Panic(args ...interface{}) {}

func (n NoopLogger) Fatal(args ...interface{}) {}

func (n NoopLogger) Debugf(template string, args ...interface{}) {}

func (n NoopLogger) Infof(template string, args ...interface{}) {}

func (n NoopLogger) Warnf(template string, args ...interface{}) {}

func (n NoopLogger) Errorf(template string, args ...interface{}) {}

func (n NoopLogger) DPanicf(template string, args ...interface{}) {}

func (n NoopLogger) Panicf(template string, args ...interface{}) {}

func (n NoopLogger) Fatalf(template string, args ...interface{}) {}

func (n NoopLogger) Debugw(msg string, keysAndValues ...interface{}) {}

func (n NoopLogger) Infow(msg string, keysAndValues ...interface{}) {}

func (n NoopLogger) Warnw(msg string, keysAndValues ...interface{}) {}

func (n NoopLogger) Errorw(msg string, keysAndValues ...interface{}) {}

func (n NoopLogger) DPanicw(msg string, keysAndValues ...interface{}) {}

func (n NoopLogger) Panicw(msg string, keysAndValues ...interface{}) {}

func (n NoopLogger) Fatalw(msg string, keysAndValues ...interface{}) {}

// used for testing
func NewNoopLogger() log.Logger {
	return NoopLogger{}
}

func SetUpConfigBeforeTestWithConfigPath(configPath, fileName string) error {
	cliCtx := context.GetContext()
	cliCtx.ConfigFile = fileName
	cliCtx.ConfigDir = configPath
	err := config.InitConfig(cliCtx)
	if err != nil {
		path := cliCtx.ConfigFilename()
		return fmt.Errorf("failed to load config file from %s", path)
	}
	return nil
}
