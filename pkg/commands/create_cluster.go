package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"github.com/RafaySystems/rcloud-cli/pkg/utils"
)

// define flags names here
const (
	CreateClusterTypeFlag = "type"
)

type clusterYamlConfig struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Labels  map[string]string `yaml:"labels"`
		Name    string            `yaml:"name"`
		Project string            `yaml:"project"`
	} `yaml:"metadata"`
	Spec struct {
		ClusterType          string `yaml:"clusterType"`
		Location             string `yaml:"location"`
		KubernetesProvider   string `yaml:"kubernetesProvider"`
		ProvisionEnvironment string `yaml:"provisionEnvironment"`
	} `yaml:"spec"`
}

type CreateClusterOptions struct {
	YamlConfigPath string
	logger         log.Logger
}

func NewCreateClusterOptions(logger log.Logger) *CreateClusterOptions {
	o := new(CreateClusterOptions)
	o.logger = logger
	return o
}

func (o *CreateClusterOptions) Validate(cmd *cobra.Command, args []string) error {
	// retrieve the flags
	flagSet := cmd.Flags()
	if flagSet.Changed(YamlConfigFlag) {
		return cobra.NoArgs(cmd, args)
	}

	// number of args must be 2
	f := cobra.ExactArgs(1)
	if err := f(cmd, args); err != nil {
		return err
	}

	if args[0] != "imported" && args[0] != "gke" && args[0] != "mks" {
		return fmt.Errorf("other cluster types except for [imported|gke] is not supported, type is %s", args[0])
	}

	return nil
}

func (o *CreateClusterOptions) Run(cmd *cobra.Command, args []string) error {
	flagSet := cmd.Flags()

	if flagSet.Changed(YamlConfigFlag) {
		return clusterConfigFile(cmd)
	}

	return nil
}

func clusterConfigFile(cmd *cobra.Command) error {
	flagSet := cmd.Flags()

	fn, err := flagSet.GetString(YamlConfigFlag)
	if err != nil {
		return err
	}
	// check if the file exists
	if !utils.FileExists(fn) {
		return fmt.Errorf("file %s does not exist", fn)
	}
	// make sure the file is a yaml file
	if filepath.Ext(fn) != ".yml" && filepath.Ext(fn) != ".yaml" {
		return fmt.Errorf("file must a yaml file, file type is %s", filepath.Ext(fn))
	}
	// read the file
	if f, err := os.Open(fn); err == nil {
		// capture the entire file
		fc, err := ioutil.ReadAll(f)
		log.GetLogger().Debugf("file is:\n%s", fc)
		if err != nil {
			return err
		}

		// unmarshal the data
		var c clusterYamlConfig
		err = yaml.Unmarshal(fc, &c)
		if err != nil {
			return err
		}

		log.GetLogger().Debugf("Unmarshalled struct: %v", c)
		if c.Kind != "Cluster" {
			return fmt.Errorf("please provide a correct yaml config for resource cluster, kind was %s", c.Kind)
		}

		if c.Spec.ClusterType == "imported" {
			f.Close()
			return importClusterConfigFile(cmd)
		}

		f.Close()
		return fmt.Errorf("other cluster types except for [imported] is not supported yet, type is %s", c.Spec.ClusterType)

	}

	return err
}

func (o *CreateClusterOptions) AddFlags(cmd *cobra.Command) {
	AddYamlConfigFlagVar(&o.YamlConfigPath, cmd, "Use this flag to create a cluster using a YAML file")
}
