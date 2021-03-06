package commands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/log"
	"github.com/spf13/cobra"
)

const (
	DownloadKubeconfigClusterFlag            = "cluster"
	DownloadKubeconfigNamespaceShorthandFlag = "n"
	DownloadKubeconfigToFileFlag             = "to-file"
	DownloadKubeconfigToFileShorthandFlag    = "s"
)

type DownloadKubeconfigOptions struct {
	Cluster,
	FileOutput string
	logger log.Logger
}

func NewDownloadKubeconfigOptions(logger log.Logger) *DownloadKubeconfigOptions {
	o := new(DownloadKubeconfigOptions)
	o.logger = logger
	return o
}

func (c *DownloadKubeconfigOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.ExactArgs(0)(cmd, args)
}

func (c *DownloadKubeconfigOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Infof("Start [%s]", cmd.CommandPath())

	auth := config.GetConfig().GetAppAuthProfile()

	cluster, _ := cmd.Flags().GetString("cluster")
	params := url.Values{}
	if cluster != "" {
		params.Add("opts.selector", fmt.Sprintf("paralus.dev/clusterName=%s", cluster))
	} else {
		return fmt.Errorf("cluster name not provided")
	}

	uri := fmt.Sprintf("/v2/sentry/kubeconfig/user?%s", params.Encode())
	resp, err := auth.AuthAndRequestFullResponse(uri, "GET", nil)
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig")
	}

	jsonData := &struct {
		Data string `json:"data"`
	}{}

	err = resp.JSON(jsonData)
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig")
	}

	decoded, err := base64.StdEncoding.DecodeString(jsonData.Data)
	if err != nil {
		return fmt.Errorf("failed to get kubeconfig")
	}
	yaml := string(decoded)

	toFile, _ := cmd.Flags().GetString("to-file")
	if len(toFile) != 0 {
		err := ioutil.WriteFile(toFile, []byte(yaml), 0644)
		if err != nil {
			return fmt.Errorf("failed to store the downloaded kubeconfig file ")
		}
		fmt.Println(fmt.Sprintf("kubeconfig downloaded to file - %s", toFile))
	} else {
		fmt.Println(yaml)
	}

	log.GetLogger().Infof("End [%s]", cmd.CommandPath())
	return nil
}

func (c *DownloadKubeconfigOptions) AddFlags(cmd *cobra.Command) {
	// add flags
	flagSet := cmd.Flags()
	flagSet.StringVar(&c.Cluster, "cluster", "", "Set the cluster to get kubeconfig for a specific cluster")
	flagSet.StringVarP(&c.FileOutput, "to-file", "s", "", "File location to download the kubeconfig to")
}
