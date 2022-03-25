package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/RafayLabs/rcloud-cli/pkg/cluster"
	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/location"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	CreateClusterImportedLocationFlag          = "location"
	CreateClusterImportedLocationShorthandFlag = "l"
	CreateClusterImportedTypeFlag              = "type"
)

type CreateClusterImportedOptions struct {
	logger log.Logger
}

func NewCreateClusterImportedOptions(logger log.Logger) *CreateClusterImportedOptions {
	o := new(CreateClusterImportedOptions)
	o.logger = logger
	return o
}

func (o *CreateClusterImportedOptions) Validate(cmd *cobra.Command, args []string) error {
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

	return nil
}

func (o *CreateClusterImportedOptions) Run(cmd *cobra.Command, args []string) error {
	flagSet := cmd.Flags()

	// retrieve the project id
	p, err := config.GetProjectIdFromFlagAndConfig(cmd)
	if err != nil {
		return err
	}

	if flagSet.Changed(YamlConfigFlag) {
		return importClusterConfigFile(cmd)
	}

	// get cluster name
	n := args[0]

	// get location
	l, _ := flagSet.GetString(CreateClusterImportedLocationFlag)

	// check if location exists
	if _, err := location.GetLocation(l); err != nil {
		return err
	}

	// get type
	cType, err := flagSet.GetString(CreateClusterImportedTypeFlag)
	if err != nil {
		cType = ""
	} else {
		if cType != "aks" && cType != "gke" && cType != "openshift" && cType != "" {
			return fmt.Errorf("imported cluster type %s is not supported", cType)
		}
	}

	// create the imported cluster
	switch cType {
	case "aks":
		_, err = cluster.NewImportClusterAKS(n, l, p)
		if err != nil {
			return err
		}
	case "gke":
		_, err = cluster.NewImportClusterGKE(n, l, p)
		if err != nil {
			return err
		}
	case "openshift":
		_, err = cluster.NewImportClusterOpenshift(n, l, p)
		if err != nil {
			return err
		}
	default:
		_, err = cluster.NewImportCluster(n, l, p)
		if err != nil {
			return err
		}
	}

	// get the bootstrap file
	bf, err := cluster.GetBootstrapFile(n, p)
	if err != nil {
		return err
	}
	fmt.Println(bf)
	return nil
}

func importClusterConfigFile(cmd *cobra.Command) error {
	flagSet := cmd.Flags()

	// retrieve the project id
	_, err := config.GetProjectIdFromFlagAndConfig(cmd)
	if err != nil {
		return err
	}

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

		if c.Spec.ClusterType != "imported" {
			return fmt.Errorf("cluster types in file is not imported, type is %s", c.Spec.ClusterType)
		}

		// check if project is provided
		if c.Metadata.Project != "" {
			_, err = config.GetProjectIdByName(c.Metadata.Project)
			if err != nil {
				return err
			}
		}

		if c.Spec.Location == "" {
			return fmt.Errorf("location is required")
		}

		// check if location exists
		if _, err := location.GetLocation(c.Spec.Location); err != nil {
			return err
		}

		// create the cluster
		_, err = cluster.NewImportCluster(c.Metadata.Name, c.Spec.Location, c.Metadata.Project)

		if err != nil {
			return err
		}

		// get the bootstrap file
		bf, err := cluster.GetBootstrapFile(c.Metadata.Name, c.Metadata.Project)
		if err != nil {
			return err
		}
		fmt.Println(bf)
		return nil
	}

	return err
}

func (o *CreateClusterImportedOptions) AddFlags(cmd *cobra.Command) {
	// define flags
	flagSet := cmd.PersistentFlags()
	flagSet.StringP(CreateClusterImportedTypeFlag, "", "", "type of imported cluster aks|gke|openshift")
	flagSet.StringP(CreateClusterImportedLocationFlag, CreateClusterImportedLocationShorthandFlag, "sanjose-us",
		"Location to set. Will cause the command to error out if the location doesn't exist")
}
