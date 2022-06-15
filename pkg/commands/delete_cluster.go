package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/paralus/cli/pkg/cluster"
	"github.com/paralus/cli/pkg/config"
	"github.com/paralus/cli/pkg/log"
	"github.com/spf13/cobra"
)

// TODO: include a query to ask for confirmation for deletion of a cluster, as well as a --approve flag

// define flags names here
const (
	ClusterDeleteConfirmFlag      = "yes"
	ClusterDeleteConfirmShortFlag = "y"
)

type DeleteClusterOptions struct {
	logger log.Logger
}

func NewDeleteClusterOptions(logger log.Logger) *DeleteClusterOptions {
	o := new(DeleteClusterOptions)
	o.logger = logger
	return o
}

func (c *DeleteClusterOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MinimumNArgs(1)(cmd, args)
}

func (c *DeleteClusterOptions) Run(cmd *cobra.Command, args []string) error {
	var y bool
	flagSet := cmd.Flags()

	// retrieve the project id
	p, err := config.GetProjectIdFromFlagAndConfig(cmd)
	if err != nil {
		return err
	}

	y, _ = flagSet.GetBool(ClusterDeleteConfirmFlag)

	if !y {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are you sure you want to perform this operation? (y/n): ")
		text, _ := reader.ReadString('\n')
		if text != "" && (text[0] == 'y' || text[0] == 'Y') {
			y = true
		}
	}

	if !y {
		return nil
	}

	// delete the specified clusters
	for _, a := range args {
		if err := cluster.DeleteCluster(a, p); err != nil {
			if !strings.Contains(err.Error(), "cluster not found") {
				return err
			}
		} else {
			log.GetLogger().Infof("Deleted %s", a)
		}
	}
	return nil
}

func (c *DeleteClusterOptions) AddFlags(cmd *cobra.Command) {
	// define flags
	flagSet := cmd.PersistentFlags()
	flagSet.BoolP(ClusterDeleteConfirmFlag, ClusterDeleteConfirmShortFlag, false,
		"Delete the cluster without confirmation")
}
