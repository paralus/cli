package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/RafaySystems/rcloud-cli/pkg/idp"
	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"github.com/RafaySystems/rcloud-cli/pkg/printer"
	"github.com/spf13/cobra"
)

// GetIdpOptions flagpole
type GetIdpOptions struct {
	Limit, Offset int
	logger        log.Logger
}

// Validate does argument validation
// this function is reserved for "offline" validation, meaning doing validations on
// the arguments without making an network calls.
func (o *GetIdpOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MaximumNArgs(1)(cmd, args)
}

func (o *GetIdpOptions) Run(cmd *cobra.Command, args []string) (err error) {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())
	flags := cmd.Flags()

	if len(args) == 1 {
		resp, err := idp.GetIdpByName(args[0])
		if err != nil {
			return fmt.Errorf("failed to get Idp API response")
		}

		ibytes, err := json.Marshal(resp)
		if err != nil {
			return fmt.Errorf("failed to get idp %s", args[0])
		}
		err = idp.Print(cmd, ibytes)
		if err != nil {
			return fmt.Errorf("failed to get idp %s", args[0])
		}

		log.GetLogger().Debugf("End [%s]", cmd.CommandPath())
		return nil

	} else {
		idps, err := idp.ListIdpWithCmd(cmd)
		if err != nil {
			return fmt.Errorf("failed to get Idps")
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
			"domain",
			"group_attribute_name",
			"acsURL",
		}

		rows := make([][]string, 0, len(idps.Items))
		// set the rows
		for _, p := range idps.Items {
			rows = append(rows, []string{
				p.Metadata.Name,
				p.Spec.Domain,
				p.Spec.GroupAttributeName,
				p.Spec.AcsUrl,
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
func (o *GetIdpOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.Flags()
	flagSet.IntVarP(&o.Limit, "limit", "l", 100, "Set the pagination limit")
	flagSet.IntVarP(&o.Offset, "offset", "e", 0, "Set the pagination offset")
}

// NewGetIdpOptions is used to create an new instance of the command.
// Parameters are dependencies required by the command, like a logger
func NewGetIdpOptions(logger log.Logger) *GetIdpOptions {
	options := new(GetIdpOptions)
	options.logger = logger
	return options
}
