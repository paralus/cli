package commands

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"syscall"

	"github.com/paralus/cli/pkg/kratos"
	"github.com/paralus/cli/pkg/log"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	downloadConfigFilePathFlag = "to-file"
)

type DownloadConfigsOptions struct {
	downloadConfigFilePath string
	logger                 log.Logger
}

func NewDownloadConfigOptions(logger log.Logger) *DownloadConfigsOptions {
	o := new(DownloadConfigsOptions)
	o.logger = logger
	return o
}

func (o *DownloadConfigsOptions) Validate(cmd *cobra.Command, args []string) error {
	f := cobra.ExactArgs(1)
	if err := f(cmd, args); err != nil {
		return err
	}
	_, err := url.ParseRequestURI(args[0])
	if err != nil {
		return err
	}

	return nil
}

func (o *DownloadConfigsOptions) Run(cmd *cobra.Command, args []string) error {
	var email string
	var password string

	fmt.Print("Requesting credentials from user.\n")
	fmt.Print("Enter Email: ")
	fmt.Scanf("%s", &email)

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return err
	}
	password = string(bytePassword)

	kc, err := kratos.Login(args[0], email, password)
	if err != nil {
		return fmt.Errorf("failed to login. You may have entered an invalid username or password or paralus host endpoint :: %v", err.Error())
	}

	o.logger.Debug("Fetching CLI config for user.")
	res, err := kc.HttpGet(fmt.Sprintf("%s/auth/v3/cli/config", args[0]))
	if err != nil {
		return fmt.Errorf("failed to download cli config : %v", err.Error())
	}

	cliConfig, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if o.downloadConfigFilePath == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		err = os.MkdirAll(fmt.Sprintf("%s/.paralus/cli/", userHomeDir), os.ModePerm)
		if err != nil {
			return err
		}
		o.downloadConfigFilePath = fmt.Sprintf("%s/.paralus/cli/config.json", userHomeDir)
	}

	file, err := os.Create(o.downloadConfigFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(cliConfig)
	if err != nil {
		return err
	}

	fmt.Printf("\nCLI config stored at `%s`\n", o.downloadConfigFilePath)

	return nil
}

func (o *DownloadConfigsOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.PersistentFlags()
	flagSet.StringVar(&o.downloadConfigFilePath, downloadConfigFilePathFlag, "",
		"File location to download CLI config")
}
