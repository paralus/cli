package cmd

import (
	"github.com/paralus/cli/pkg/commands"
	"github.com/spf13/cobra"
)

func newUpdateGroupassociationCmd(o commands.CmdOptions) *cobra.Command {
	// cmd represents the update command
	cmd := &cobra.Command{
		Use:     "groupassociation <group-name>",
		Aliases: []string{"ga"},
		Short:   "Update a group association",
		Long:    "Update a group association",
		Example: `
Using command to associate project:
	pctl update groupassociation sample-group --associateproject sample-proj --roles ADMIN
	pctl update groupassociation sample-group --associateproject sample-proj --roles PROJECT_READ_ONLY,INFRA_ADMIN 
	pctl update groupassociation sample-group --associateproject sample-proj --roles NAMESPACE_READ_ONLY --namespace samplens

roles:
	ADMIN
	PROJECT_ADMIN
	PROJECT_READ_ONLY
	INFRA_ADMIN
	INFRA_READ_ONLY
	NAMESPACE_READ_ONLY
	NAMESPACE_ADMIN

Using command to associate user:
	pctl update groupassociation sample-group  --associateuser y --addusers example.user@company.co
	pctl update groupassociation sample-group  --associateuser y --addusers example.user@company.co,example.user-two@company.co --removeusers example.user-three@company.co 


`,
		Args: o.Validate,
		RunE: o.Run,
	}

	o.AddFlags(cmd)

	return cmd

}
