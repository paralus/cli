package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/rafaylabs/rcloud-cli/pkg/log"
	"github.com/rafaylabs/rcloud-cli/pkg/printer"
	"github.com/rafaylabs/rcloud-cli/pkg/project"
	"github.com/spf13/cobra"
)

// flagpole
type GetProjectOptions struct {
	Limit, Offset int
	logger        log.Logger
}

// Validate does argument validation
// this function is reserved for "offline" validation, meaning doing validations on
// the arguments without making an network calls.
func (o *GetProjectOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MaximumNArgs(1)(cmd, args)
}

func (o *GetProjectOptions) Run(cmd *cobra.Command, args []string) (err error) {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())
	flags := cmd.Flags()

	if len(args) == 1 {
		projectName := args[0]
		proj, err := project.GetProjectByName(projectName)
		if err != nil {
			return fmt.Errorf("failed to get project %s", projectName)
		}

		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, proj, cmd.OutOrStdout())
			return nil
		}

		pbytes, err := json.Marshal(proj)
		if err != nil {
			return fmt.Errorf("failed to get project %s", projectName)
		}
		err = project.Print(cmd, []byte(pbytes))
		if err != nil {
			return fmt.Errorf("failed to get project %s", projectName)
		}
	} else {
		projects, err := project.ListProjectsWithCmd(cmd)
		if err != nil {
			return fmt.Errorf("failed to get projects %s", err.Error())
		}
		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, projects, cmd.OutOrStdout())
			return nil
		}

		// set the columns
		c := []string{
			"project",
			"default",
		}

		rows := make([][]string, 0, len(projects.Items))
		// set the rows
		for _, p := range projects.Items {
			rows = append(rows, []string{
				p.Metadata.Name,
				strconv.FormatBool(p.Spec.Default),
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

// AddFlags is where you define the command flags and attach the flagpole
func (o *GetProjectOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.Flags()
	flagSet.IntVarP(&o.Limit, "limit", "l", 100, "Set the pagenation limit")
	flagSet.IntVarP(&o.Offset, "offset", "e", 0, "Set the pagenation offset")
}

// NewGetProjectOptions is used to create an new instance of the command.
// Parameters are dependencies required by the command, like a logger
func NewGetProjectOptions(logger log.Logger) CmdOptions {
	options := new(GetProjectOptions)
	options.logger = logger
	return options
}
