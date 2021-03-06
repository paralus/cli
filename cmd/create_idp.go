package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newCreateIdpCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the create command
	cmd := &cobra.Command{
		Use:   "idp <idp-name>",
		Short: "Create a new idp config",
		Long:  "pctl create idp <idp-name> <domain-name> <group-name>",
		Example: `

Basic Info : This command is used to configure Identity providers (IDP)
----------------------
idp-name : Configure unique IDP name 
idp-type : 3rd party SSO configuration profiles

	Supported idp-types
	--	okta | ping | custom

domain-name: email domain of the organization e.g company.com
group-name : Set the name of the Group Attribute Statement to the group with assigned roles in the console
----------------------

Examples:
	Basic Command :
		pctl create idp <idp-name> <idp-type> <domain-name> <group-name>
		pctl create idp paralus custom paralus.com paralussuperadmins

`,

		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
