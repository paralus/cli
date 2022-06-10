package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/paralus/cli/pkg/group"
	"github.com/paralus/cli/pkg/log"
	"github.com/paralus/cli/pkg/printer"
	"github.com/spf13/cobra"
)

// flagpole
type GetGroupOptions struct {
	Limit, Offset int
	logger        log.Logger
}

// Validate does argument validation
// this function is reserved for "offline" validation, meaning doing validations on
// the arguments without making an network calls.
func (o *GetGroupOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MaximumNArgs(1)(cmd, args)
}

func (o *GetGroupOptions) Run(cmd *cobra.Command, args []string) (err error) {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())
	flags := cmd.Flags()

	if len(args) == 1 {
		groupName := args[0]
		grp, err := group.GetGroupByName(groupName)
		if err != nil {
			return fmt.Errorf("failed to get group %s", groupName)
		}

		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, grp, cmd.OutOrStdout())
			return nil
		}

		gbytes, err := json.Marshal(grp)
		if err != nil {
			return err
		}
		err = group.Print(cmd, gbytes)
		if err != nil {
			return fmt.Errorf("failed to get group %s", groupName)
		}

	} else {
		groups, err := group.ListGroupsWithCmd(cmd)
		if err != nil {
			return fmt.Errorf("failed to get groups")
		}
		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, groups, cmd.OutOrStdout())
			return nil
		}

		// set the columns
		c := []string{
			"name",
			"description",
			"type",
			"users",
		}

		// set the columns
		col := []string{
			"group",
			"project",
			"role",
		}

		rows := make([][]string, 0, len(groups.Items))
		prows := make([][]string, 0, len(groups.Items))
		// set the rows
		for _, p := range groups.Items {
			prows = make([][]string, 0, len(groups.Items))
			rows = append(rows, []string{
				p.Metadata.Name,
				p.Metadata.Description,
				p.Spec.Type,
				strings.Join(p.Spec.Users, ","),
			})

			for _, r := range p.Spec.ProjectNamespaceRoles {
				prows = append(prows, []string{
					p.Metadata.Name,
					r.GetProject(),
					r.GetRole(),
				})
			}
		}
		sort.Slice(rows, func(i, j int) bool {
			return strings.Compare(rows[i][0], rows[j][0]) == -1
		})

		printer.PrintTable(c, rows, cmd.OutOrStdout())
		printer.PrintTable(col, prows, cmd.OutOrStdout())

	}
	log.GetLogger().Debugf("End [%s]", cmd.CommandPath())
	return nil
}

// AddFlags is where you define the command flags and attach the flagpole
func (o *GetGroupOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.Flags()
	flagSet.IntVarP(&o.Limit, "limit", "l", 100, "Set the pagenation limit")
	flagSet.IntVarP(&o.Offset, "offset", "e", 0, "Set the pagenation offset")
}

// NewGetGroupOptions is used to create an new instance of the command.
// Parameters are dependencies required by the command, like a logger
func NewGetGroupOptions(logger log.Logger) *GetGroupOptions {
	options := new(GetGroupOptions)
	options.logger = logger
	return options
}
