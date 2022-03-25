package cmd

import (
	"github.com/RafayLabs/rcloud-cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newCreateOIDCProviderCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the create command
	cmd := &cobra.Command{
		Use:   "oidc <oidc-provider-name>",
		Short: "Create a new oidc provider config",
		Long:  "rctl create oidc <oidc-provider-name> --clientid <client-id> --callback-url <callback-url> --scopes <scopes>",
		Example: `

Basic Info : This command is used to configure OIDC Identity providers (IDP)
----------------------
oidc-provider-name : Configure unique IDP name 

client-id: client identifier which was created while registering in the provider
callback-url : Set the callback url
scopes: scopes that are required
----------------------

Examples:
	Basic Command :
		rctl create oidc github --clientid 721396hsad8721wjhad8 --callback-url http://myownweburl.com/cb --scopes name
`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
