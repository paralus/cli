package commands

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/rafaylabs/rcloud-cli/pkg/log"
	"github.com/rafaylabs/rcloud-cli/pkg/printer"
	"github.com/rafaylabs/rcloud-cli/pkg/user"
	"github.com/spf13/cobra"
)

// flagpole
type GetUserOptions struct {
	Limit, Offset int
	logger        log.Logger
}

// Validate does argument validation
// this function is reserved for "offline" validation, meaning doing validations on
// the arguments without making an network calls.
func (o *GetUserOptions) Validate(cmd *cobra.Command, args []string) error {
	return cobra.MaximumNArgs(1)(cmd, args)
}

func (o *GetUserOptions) Run(cmd *cobra.Command, args []string) (err error) {
	log.GetLogger().Debugf("Start [%s]", cmd.CommandPath())
	flags := cmd.Flags()

	if len(args) == 1 {
		userName := args[0]
		usr, err := user.GetUserByName(userName)
		if err != nil {
			return fmt.Errorf("failed to get user %s", userName)
		}

		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, usr, cmd.OutOrStdout())
			return nil
		}

		ubytes, err := json.Marshal(usr)
		if err != nil {
			return fmt.Errorf("failed to get user %s", userName)
		}
		err = user.Print(cmd, ubytes)
		if err != nil {
			return fmt.Errorf("failed to get user %s", userName)
		}

	} else {
		users, err := user.ListUsersWithCmd(cmd)
		if err != nil {
			return fmt.Errorf("failed to get users")
		}
		// json or yaml output
		if flags.Changed("output") {
			o, err := flags.GetString("output")
			if err != nil {
				return err
			}
			printer.PrintOutputJsonOrYaml(o, users, cmd.OutOrStdout())
			return nil
		}

		// set the columns
		c := []string{
			"name",
			"first_name",
			"last_name",
			"groups",
		}

		rows := make([][]string, 0, len(users.Items))
		// set the rows
		for _, p := range users.Items {
			rows = append(rows, []string{
				p.Metadata.Name,
				p.Spec.FirstName,
				p.Spec.LastName,
				strings.Join(p.Spec.Groups, ","),
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
func (o *GetUserOptions) AddFlags(cmd *cobra.Command) {
	flagSet := cmd.Flags()
	flagSet.IntVarP(&o.Limit, "limit", "l", 100, "Set the pagenation limit")
	flagSet.IntVarP(&o.Offset, "offset", "e", 0, "Set the pagenation offset")
}

// NewGetUserOptions is used to create an new instance of the command.
// Parameters are dependencies required by the command, like a logger
func NewGetUserOptions(logger log.Logger) *GetUserOptions {
	options := new(GetUserOptions)
	options.logger = logger
	return options
}
