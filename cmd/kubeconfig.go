package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"

	"github.com/RafaySystems/rcloud-cli/pkg/config"
	"github.com/RafaySystems/rcloud-cli/pkg/exit"
	"github.com/RafaySystems/rcloud-cli/pkg/log"

	"github.com/spf13/cobra"
)

// KubeconfigCmd is command definition for kubeconfig root command
var KubeconfigCmd = &cobra.Command{
	Use:     "kubeconfig",
	Short:   "Generate kubeconfig",
	Long:    "Allows the user to generate a kubeconfig",
	Aliases: []string{"kc"},
}

// KubeconfigDownloadCmd is command definition for kubeconfig download cmd
var KubeconfigDownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the generated kubeconfig",
	Long:  "Download the generated kubeconfig",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		log.GetLogger().Infof("Start [%s]", cmd.CommandPath())

		auth := config.GetConfig().GetAppAuthProfile()

		defaultNamespace, _ := cmd.Flags().GetString("namespace")
		cluster, _ := cmd.Flags().GetString("cluster")
		params := url.Values{}
		if defaultNamespace != "" {
			params.Add("namespace", defaultNamespace)
		}
		if cluster != "" {
			params.Add("opts.selector", fmt.Sprintf("rafay.dev/clusterName=%s", cluster))
			params.Add("opts.ID", config.GetConfig().APIKey)
			params.Add("opts.organization", config.GetConfig().Organization)
		}

		uri := fmt.Sprintf("/v2/sentry/kubeconfig/user?%s", params.Encode())
		resp, err := auth.AuthAndRequestFullResponse(uri, "GET", nil)
		if err != nil {
			exit.SetExitWithError(err, "failed to get kubeconfig")
			return
		}

		jsonData := &struct {
			Data string `json:"data"`
		}{}

		err = resp.JSON(jsonData)
		if err != nil {
			exit.SetExitWithError(err, "failed to get kubeconfig")
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(jsonData.Data)
		if err != nil {
			exit.SetExitWithError(err, "failed to get kubeconfig")
			return
		}
		yaml := string(decoded)

		toFile, _ := cmd.Flags().GetString("to-file")
		if len(toFile) != 0 {
			err := ioutil.WriteFile(toFile, []byte(yaml), 0644)
			if err != nil {
				exit.SetExitWithError(err, "failed to store the downloaded kubeconfig file ")
				return
			}
			fmt.Println(fmt.Sprintf("kubeconfig downloaded to file - %s", toFile))
		} else {
			fmt.Println(yaml)
		}

		log.GetLogger().Infof("End [%s]", cmd.CommandPath())
	},
}

func init() {
	KubeconfigDownloadCmd.Flags().StringP("namespace", "n", "", "Set the default namespace for the kubeconfig")
	KubeconfigDownloadCmd.Flags().String("cluster", "", "Set the cluster to get kubeconfig for a specific cluster")
	KubeconfigDownloadCmd.Flags().StringP("to-file", "l", "", "File location to download the kubeconfig to")
	KubeconfigCmd.AddCommand(
		KubeconfigDownloadCmd,
	)
}
