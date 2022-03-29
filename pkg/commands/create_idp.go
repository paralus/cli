package commands

import (
	"fmt"
	"strings"

	commonv3 "github.com/RafayLabs/rcloud-base/proto/types/commonpb/v3"
	systemv3 "github.com/RafayLabs/rcloud-base/proto/types/systempb/v3"
	"github.com/RafayLabs/rcloud-cli/pkg/config"
	"github.com/RafayLabs/rcloud-cli/pkg/idp"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/spf13/cobra"
)

type CreateIDpOptions struct {
	logger log.Logger
	config *config.Config
}

func NewCreateIdpOptions(logger log.Logger, config *config.Config) *CreateIDpOptions {
	o := new(CreateIDpOptions)
	o.logger = logger
	o.config = config
	return o
}

func (c *CreateIDpOptions) AddFlags(_ *cobra.Command) {}

func (o *CreateIDpOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MinimumNArgs(4)(cmd, args)

}

func (o *CreateIDpOptions) Run(cmd *cobra.Command, args []string) error {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())

	//.com or .xyz domain should exist
	if !strings.Contains(args[2], ".") {
		return fmt.Errorf("domain name missing in %s", args[2])
	}

	name := args[0]
	domain := args[2]
	groupname := args[3]

	idprovider := &systemv3.Idp{
		Kind: "Idp",
		Metadata: &commonv3.Metadata{
			Name:         name,
			Organization: config.GetConfig().Organization,
		},
		Spec: &systemv3.IdpSpec{
			IdpName:            name,
			Domain:             domain,
			GroupAttributeName: groupname,
		},
	}

	err := idp.CreateIdp(idprovider)
	if err != nil {
		return fmt.Errorf("failed to create idp, cause: %s", err.Error())
	}
	return nil
}
