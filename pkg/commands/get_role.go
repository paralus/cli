package commands

import (
	"sort"
	"strconv"
	"strings"

	"github.com/RafaySystems/rcloud-cli/pkg/log"
	"github.com/RafaySystems/rcloud-cli/pkg/printer"
	"github.com/RafaySystems/rcloud-cli/pkg/role"
	"github.com/spf13/cobra"
)

type GetRolesOptions struct {
	logger log.Logger
}

func NewGetRolesOptions(logger log.Logger) *GetRolesOptions {
	o := new(GetRolesOptions)
	o.logger = logger
	return o
}

func (o *GetRolesOptions) Validate(_ *cobra.Command, _ []string) error {
	return nil
}

func (o *GetRolesOptions) Run(cmd *cobra.Command, args []string) error {
	// get flags
	flags := cmd.Flags()

	// list of roles
	roles, err := role.ListRolesWithCmd(cmd)
	if err != nil {
		return err
	}

	// json or yaml output
	if flags.Changed("output") {
		o, err := flags.GetString("output")
		if err != nil {
			return err
		}
		printer.PrintOutputJsonOrYaml(o, roles, cmd.OutOrStdout())
		return nil
	}

	// set the columns
	c := []string{
		"name",
		"description",
		"is_global",
		"scope",
	}

	rows := make([][]string, 0, len(roles.Items))
	// set the rows
	for _, c := range roles.Items {
		rows = append(rows, []string{
			c.Metadata.Name,
			c.Metadata.Description,
			strconv.FormatBool(c.Spec.IsGlobal),
			c.Spec.Scope,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return strings.Compare(rows[i][0], rows[j][0]) == -1
	})

	printer.PrintTable(c, rows, cmd.OutOrStdout())

	return nil
}

func (o *GetRolesOptions) AddFlags(_ *cobra.Command) {}
