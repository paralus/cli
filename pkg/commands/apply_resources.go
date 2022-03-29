package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	infrav3 "github.com/RafayLabs/rcloud-base/proto/types/infrapb/v3"
	rolev3 "github.com/RafayLabs/rcloud-base/proto/types/rolepb/v3"
	systemv3 "github.com/RafayLabs/rcloud-base/proto/types/systempb/v3"
	userv3 "github.com/RafayLabs/rcloud-base/proto/types/userpb/v3"
	"github.com/RafayLabs/rcloud-cli/pkg/cluster"
	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/constants"
	"github.com/RafayLabs/rcloud-cli/pkg/group"
	"github.com/RafayLabs/rcloud-cli/pkg/idp"
	"github.com/RafayLabs/rcloud-cli/pkg/location"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/oidc"
	"github.com/RafayLabs/rcloud-cli/pkg/project"
	"github.com/RafayLabs/rcloud-cli/pkg/role"
	"github.com/RafayLabs/rcloud-cli/pkg/user"
	"github.com/RafayLabs/rcloud-cli/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ApplyResourcesOptions struct {
	logger         log.Logger
	config         config.Config
	YamlConfigPath string
}

func NewApplyResourcesOptions(logger log.Logger, config *config.Config) *ApplyResourcesOptions {
	o := new(ApplyResourcesOptions)
	o.logger = logger
	o.config = *config
	return o
}

func (o *ApplyResourcesOptions) Validate(cmd *cobra.Command, args []string) error {
	// retrieve the flags
	flagSet := cmd.Flags()

	if flagSet.Changed(constants.FILE_FLAG_NAME) {
		return cobra.NoArgs(cmd, args)
	}

	// number of args must be 2
	f := cobra.ExactArgs(1)
	if err := f(cmd, args); err != nil {
		return err
	}

	return nil
}

func (o *ApplyResourcesOptions) Run(cmd *cobra.Command, args []string) error {
	flagSet := cmd.Flags()

	if flagSet.Changed(constants.FILE_FLAG_NAME) {
		return processConfigFile(cmd)
	} else {
		return fmt.Errorf("invalid command! example usage: rctl apply -f file.yml")
	}
}

func processConfigFile(cmd *cobra.Command) error {

	bytes, err := readConfigFile(cmd)
	if err != nil {
		fmt.Errorf("error reading config file, cause: %s", err.Error())
	}

	// unmarshal the data
	var gvk gvkYamlConfig
	err = yaml.Unmarshal(bytes, &gvk)
	if err != nil {
		return fmt.Errorf("error reading kind, cause: %s", err.Error())
	}
	log.GetLogger().Debugf("kind is:\n%s", gvk.Kind)
	if StringInSlice(gvk.Kind, resources) {
		//process clusters
		if gvk.Kind == "Cluster" {
			var clstr infrav3.Cluster
			err = yaml.Unmarshal(bytes, &clstr)
			if err != nil {
				fmt.Errorf("error reading config file, cause: %s", err.Error())
			}

			//validate data
			if clstr.Spec.ClusterType != "imported" {
				return fmt.Errorf("cluster types in file is not imported, type is %s", clstr.Spec.ClusterType)
			}

			// check if project is provided
			if clstr.Metadata.Project != "" {
				_, err = project.GetProjectByName(clstr.Metadata.Project)
				if err != nil {
					return err
				}
			}

			if clstr.Spec.Metro != nil {
				// check if location exists
				if _, err := location.GetLocation(clstr.Spec.Metro.Name); err != nil {
					return err
				}
			}

			existing, _ := cluster.GetCluster(clstr.Metadata.Name, clstr.Metadata.Project)
			if existing != nil {
				log.GetLogger().Debugf("updating cluster: %s", clstr.Metadata.Name)
				err := cluster.UpdateCluster(&clstr)
				if err != nil {
					fmt.Printf("Error configuring resource %s due to %s \n", gvk.Kind, err.Error())
				} else {
					fmt.Printf("Resource %s of type %s configured.\n", clstr.Metadata.Name, gvk.Kind)
				}

			} else {
				log.GetLogger().Debugf("creating cluster: %s", clstr.Metadata.Name)
				err := cluster.CreateCluster(&clstr)
				if err != nil {
					fmt.Printf("Error configuring resource %s due to %s \n", gvk.Kind, err.Error())
				} else {
					fmt.Printf("Resource %s of type %s configured.\n", clstr.Metadata.Name, gvk.Kind)
				}
			}

		} else if gvk.Kind == "Project" {

			var proj systemv3.Project
			err = yaml.Unmarshal(bytes, &proj)
			if err != nil {
				fmt.Errorf("error reading config file, cause: %s", err.Error())
			}
			err := project.ApplyProject(&proj)
			if err != nil {
				fmt.Printf("Error configuring resource %s due to %s \n", gvk.Kind, err.Error())
			} else {
				fmt.Printf("Resource %s of type %s configured.\n", proj.Metadata.Name, gvk.Kind)
			}

		} else if gvk.Kind == "User" {

			var usr userv3.User
			err = yaml.Unmarshal(bytes, &usr)
			if err != nil {
				fmt.Errorf("error reading config file, cause: %s", err.Error())
			}
			err := user.ApplyUser(&usr)
			if err != nil {
				fmt.Printf("Error configuring resource %s due to %s \n", gvk.Kind, err.Error())
			} else {
				fmt.Printf("Resource %s of type %s configured.\n", usr.Metadata.Name, gvk.Kind)
			}

		} else if gvk.Kind == "Group" {

			var grp userv3.Group
			err = yaml.Unmarshal(bytes, &grp)
			if err != nil {
				fmt.Errorf("error reading config file, cause: %s", err.Error())
			}
			err := group.ApplyGroup(&grp)
			if err != nil {
				fmt.Printf("Error configuring resource %s due to %s \n", gvk.Kind, err.Error())
			} else {
				fmt.Printf("Resource %s of type %s configured.\n", grp.Metadata.Name, gvk.Kind)
			}

		} else if gvk.Kind == "Idp" {

			var id systemv3.Idp
			err = yaml.Unmarshal(bytes, &id)
			if err != nil {
				fmt.Errorf("error reading config file, cause: %s", err.Error())
			}
			err := idp.ApplyIDP(&id)
			if err != nil {
				fmt.Printf("Error configuring resource %s due to %s \n", gvk.Kind, err.Error())
			} else {
				fmt.Printf("Resource %s of type %s configured.\n", id.Metadata.Name, gvk.Kind)
			}

		} else if gvk.Kind == "OIDCProvider" {
			var oidcp systemv3.OIDCProvider
			err = yaml.Unmarshal(bytes, &oidcp)
			if err != nil {
				fmt.Errorf("error reading config file, cause: %s", err.Error())
			}
			err := oidc.ApplyOIDC(&oidcp)
			if err != nil {
				fmt.Printf("Error configuring resource %s error: %s \n", gvk.Kind, err.Error())
			} else {
				fmt.Printf("Resource %s of type %s configured.\n", oidcp.Metadata.Name, gvk.Kind)
			}

		} else if gvk.Kind == "Role" {
			var r rolev3.Role
			err = yaml.Unmarshal(bytes, &r)
			if err != nil {
				fmt.Errorf("error reading config file, error: %s", err.Error())
			}
			err := role.ApplyRole(&r)
			if err != nil {
				fmt.Printf("Error configuring resource %s due to %s \n", gvk.Kind, err.Error())
			} else {
				fmt.Printf("Resource %s of type %s configured.\n", r.Metadata.Name, gvk.Kind)
			}

		}

	} else {
		fmt.Errorf("unsupported resource kind %s", gvk.Kind)
	}
	return err
}

func readConfigFile(cmd *cobra.Command) ([]byte, error) {
	flagSet := cmd.Flags()
	fn, err := flagSet.GetString(constants.FILE_FLAG_NAME)
	if err != nil {
		return nil, err
	}
	// check if the file exists
	if !utils.FileExists(fn) {
		return nil, fmt.Errorf("file %s does not exist", fn)
	}
	// make sure the file is a yaml file
	if filepath.Ext(fn) != ".yml" && filepath.Ext(fn) != ".yaml" {
		return nil, fmt.Errorf("file must a yaml file, file type is %s", filepath.Ext(fn))
	}
	var fc []byte
	// read the file
	if f, err := os.Open(fn); err == nil {
		// capture the entire file
		fc, err = ioutil.ReadAll(f)
		log.GetLogger().Debugf("file is:\n%s", fc)
		if err != nil {
			return nil, err
		}
	}
	return fc, nil
}

func (c *ApplyResourcesOptions) AddFlags(_ *cobra.Command) {}
