package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/rafaylabs/rcloud-cli/pkg/log"
	"github.com/rafaylabs/rcloud-cli/pkg/printer"
	"github.com/rafaylabs/rcloud-cli/pkg/role"
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
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())
	// get flags
	flags := cmd.Flags()

	if len(args) == 1 {
		roleName := args[0]
		r, err := role.GetRoleByName(roleName)
		if err != nil {
			return fmt.Errorf("failed to get role %s", roleName)
		}

		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, r, cmd.OutOrStdout())
			return nil
		}

		gbytes, err := json.Marshal(r)
		if err != nil {
			return err
		}
		err = role.Print(cmd, gbytes)
		if err != nil {
			return fmt.Errorf("failed to get role %s", roleName)
		}
	} else {
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
	}

	log.GetLogger().Debugf("End [%s]", cmd.CommandPath())
	return nil
}

func (o *GetRolesOptions) AddFlags(_ *cobra.Command) {}
