package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/paralus/cli/pkg/log"
	"github.com/paralus/cli/pkg/oidc"
	"github.com/paralus/cli/pkg/printer"
	"github.com/spf13/cobra"
)

// GetOIDCProviderOptions flagpole
type OIDCProviderOptions struct {
	Limit, Offset int
	logger        log.Logger
}

// Validate does argument validation
// this function is reserved for "offline" validation, meaning doing validations on
// the arguments without making an network calls.
func (o *OIDCProviderOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MaximumNArgs(1)(cmd, args)
}

func (o *OIDCProviderOptions) Run(cmd *cobra.Command, args []string) (err error) {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())
	flags := cmd.Flags()

	if len(args) == 1 {
		resp, err := oidc.GetOIDCProviderByName(args[0])
		if err != nil {
			return fmt.Errorf("failed to get OIDC Provider API response")
		}

		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, resp, cmd.OutOrStdout())
			return nil
		}

		ibytes, err := json.Marshal(resp)
		if err != nil {
			return fmt.Errorf("failed to get oidc %s", args[0])
		}
		err = oidc.Print(cmd, ibytes)
		if err != nil {
			return fmt.Errorf("failed to get oidc provider %s", args[0])
		}

		log.GetLogger().Debugf("End [%s]", cmd.CommandPath())
		return nil

	} else {
		idps, err := oidc.ListOIDCWithCmd(cmd)
		if err != nil {
			return fmt.Errorf("failed to get oidc provider")
		}
		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, idps, cmd.OutOrStdout())
			return nil
		}

		// set the columns
		c := []string{
			"name",
			"client_id",
			"callback_url",
			"auth_url",
			"mapper_url",
			"issuer_url",
			"scopes",
			"requested_claims",
		}

		rows := make([][]string, 0, len(idps.Items))
		// set the rows
		for _, p := range idps.Items {
			rows = append(rows, []string{
				p.Metadata.Name,
				p.Spec.ClientId,
				p.Spec.CallbackUrl,
				p.Spec.AuthUrl,
				p.Spec.MapperUrl,
				p.Spec.IssuerUrl,
				strings.Join(p.Spec.Scopes, ","),
				p.Spec.RequestedClaims.String(),
			})
		}
		sort.Slice(rows, func(i, j int) bool {
			return strings.Compare(rows[i][0], rows[j][0]) == -1
		})

		printer.PrintTable(c, rows, cmd.OutOrStdout())
		log.GetLogger().Debugf("End [%s]", cmd.CommandPath())
		return nil
	}

}

// AddFlags is where you define the command flags and attach the flagpole
func (o *OIDCProviderOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.Flags()
	flagSet.IntVarP(&o.Limit, "limit", "l", 100, "Set the pagination limit")
	flagSet.IntVarP(&o.Offset, "offset", "e", 0, "Set the pagination offset")
}

// NewGetOIDCProviderOptions is used to create an new instance of the command.
// Parameters are dependencies required by the command, like a logger
func NewGetOIDCProviderOptions(logger log.Logger) *OIDCProviderOptions {
	options := new(OIDCProviderOptions)
	options.logger = logger
	return options
}
