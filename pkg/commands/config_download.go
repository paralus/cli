package commands

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"syscall"

	"github.com/paralus/cli/pkg/kratos"
	"github.com/paralus/cli/pkg/log"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// defaultConfigDirPath is the default directory to download config to.
var defaultConfigDirPath = path.Join(".paralus", "cli")

const (
	downloadConfigFilePathFlag = "to-file"
	emailFlag                  = "email"
	passwordFlag               = "password"
	defaultConfigFileName      = "config.json"
)

type DownloadConfigsOptions struct {
	downloadConfigFilePath string
	email                  string
	password               string
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
		return errors.New("invalid console URL.")
	}

	if fi, err := os.Stat(o.downloadConfigFilePath); err == nil && fi.IsDir() {
		return fmt.Errorf("%s is a directory", o.downloadConfigFilePath)
	}

	if o.email == "" && o.password != "" || o.email != "" && o.password == "" {
		return errors.New("please provide both email and password to login")
	}

	return nil
}

func (o *DownloadConfigsOptions) Run(cmd *cobra.Command, args []string) error {

	if o.email == "" && o.password == "" {
		fmt.Print("Enter Email: ")
		fmt.Scanf("%s", &o.email)

		fmt.Print("Enter Password: ")
		bytePassword, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return err
		}
		fmt.Println() // add newline after password
		o.password = string(bytePassword)
	}

	kc, err := kratos.Login(args[0], o.email, o.password)
	if err != nil {
		if errors.Is(err, kratos.ErrKratosFlow) {
			return fmt.Errorf("failed to login. Please verify console host")
		}
		return fmt.Errorf("failed to login. Reason: %s", err.Error())
	}

	o.logger.Info("Fetching CLI config for user.")
	res, err := kc.HttpGet(fmt.Sprintf("%s/auth/v3/cli/config", args[0]))
	if err != nil || (res.StatusCode != http.StatusOK) {
		o.logger.Infof("Error while download cli config : %v ", res.StatusCode)
		return errors.New("failed to download cli config.")
	}

	defer res.Body.Close()

	cliConfig, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if o.downloadConfigFilePath == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		downloadDirPath := path.Join(userHomeDir, defaultConfigDirPath)
		err = os.MkdirAll(downloadDirPath, os.ModePerm)
		if err != nil {
			return err
		}
		o.downloadConfigFilePath = path.Join(downloadDirPath, defaultConfigFileName)
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

	fmt.Printf("CLI config stored at `%s`\n", o.downloadConfigFilePath)
	return nil
}

func (o *DownloadConfigsOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.PersistentFlags()
	flagSet.StringVar(&o.downloadConfigFilePath, downloadConfigFilePathFlag, "",
		"File location to download CLI config")
	flagSet.StringVar(&o.email, emailFlag, "",
		"Email for login to Paralus")
	flagSet.StringVar(&o.password, passwordFlag, "",
		"Password for login to Paralus")
}
