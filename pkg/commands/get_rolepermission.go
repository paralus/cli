package commands

import (
	"sort"
	"strings"

	"github.com/paralus/cli/pkg/log"
	"github.com/paralus/cli/pkg/printer"
	"github.com/paralus/cli/pkg/rolepermission"
	"github.com/spf13/cobra"
)

type GetRolePermissionOptions struct {
	logger log.Logger
}

func NewGetRolePermissionOptions(logger log.Logger) *GetRolePermissionOptions {
	o := new(GetRolePermissionOptions)
	o.logger = logger
	return o
}

func (o *GetRolePermissionOptions) Validate(_ *cobra.Command, _ []string) error {
	return nil
}

func (o *GetRolePermissionOptions) Run(cmd *cobra.Command, args []string) error {
	// get flags
	flags := cmd.Flags()

	// list of roles
	rolepermissions, err := rolepermission.ListRolePermissionWithCmd(cmd)
	if err != nil {
		return err
	}

	// json or yaml output
	if flags.Changed("output") {
		o, err := flags.GetString("output")
		if err != nil {
			return err
		}
		printer.PrintOutputJsonOrYaml(o, rolepermissions, cmd.OutOrStdout())
		return nil
	}

	// set the columns
	c := []string{
		"name",
		"description",
	}

	rows := make([][]string, 0, len(rolepermissions.Items))
	// set the rows
	for _, c := range rolepermissions.Items {
		rows = append(rows, []string{
			c.Metadata.Name,
			c.Metadata.Description,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		return strings.Compare(rows[i][0], rows[j][0]) == -1
	})

	printer.PrintTable(c, rows, cmd.OutOrStdout())

	return nil
}

func (o *GetRolePermissionOptions) AddFlags(_ *cobra.Command) {}
