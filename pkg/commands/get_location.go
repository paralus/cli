package commands

import (
	"sort"
	"strings"

	"github.com/RafayLabs/rcloud-cli/pkg/location"
	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/printer"
	"github.com/spf13/cobra"
)

type GetLocationOptions struct {
	logger log.Logger
}

func NewGetLocationOptions(logger log.Logger) *GetLocationOptions {
	o := new(GetLocationOptions)
	o.logger = logger
	return o
}

func (o *GetLocationOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.NoArgs(cmd, args)
}

func (o *GetLocationOptions) Run(cmd *cobra.Command, args []string) error {
	// get flags
	flags := cmd.Flags()

	// get list of locations
	locations, err := location.ListAllLocation()
	if err != nil {
		return err
	}

	// json or yaml output
	if flags.Changed("output") {
		o, err := flags.GetString("output")
		if err != nil {
			return err
		}
		printer.PrintOutputJsonOrYaml(o, locations, cmd.OutOrStdout())
		return nil
	}

	// set the columns
	c := []string{
		"name",
		"city",
		"state",
		"state_code",
		"country",
		"country_code",
	}

	rows := make([][]string, 0, len(locations))
	// set the rows
	for _, l := range locations {
		rows = append(rows, []string{
			l.Name,
			l.City,
			l.State,
			l.StateCode,
			l.Country,
			l.CountryCode,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return strings.Compare(rows[i][0], rows[j][0]) == -1
	})

	printer.PrintTable(c, rows, cmd.OutOrStdout())
	return nil
}

func (o *GetLocationOptions) AddFlags(_ *cobra.Command) {}
