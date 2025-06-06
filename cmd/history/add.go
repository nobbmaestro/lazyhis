package history

import (
	"fmt"
	"os"

	appopts "github.com/nobbmaestro/lazyhis/pkg/app"

	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/spf13/cobra"
)

var historyAddOpts = &HistoryOptions{}

var historyAddCmd = &cobra.Command{
	Use:   "add [CMD...]",
	Short: "Add history record",
	Args:  cobra.ArbitraryArgs,
	RunE:  runHistoryAdd,
}

func runHistoryAdd(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	app := reg.GetApp()

	if historyAddOpts.path == "" {
		if currentPath, err := os.Getwd(); err == nil {
			historyAddOpts.path = currentPath
		}
	}

	record, err := app.AddHistory(
		false, // dryRun
		false, // verbose
		false, // addUniqueOnly
		appopts.WithQuery(args),
		appopts.WithPath(historyAddOpts.path),
		appopts.WithSession(historyAddOpts.session),
		appopts.WithExitCode(historyAddOpts.exitCode),
		appopts.WithExecutedIn(historyAddOpts.executedIn),
	)
	if err != nil {
		return err
	}
	if record != nil {
		fmt.Println(record.ID)
	}

	return nil
}

func init() {
	historyAddCmd.
		Flags().
		IntVarP(&historyAddOpts.exitCode, "exit-code", "e", -1, "exit code for the command")
	historyAddCmd.
		Flags().
		IntVarP(&historyAddOpts.executedIn, "duration", "d", -1, "execution duration of the CMD in milliseconds")
	historyAddCmd.
		Flags().
		StringVarP(&historyAddOpts.path, "path", "p", "", "working directory context")
	historyAddCmd.
		Flags().
		StringVarP(&historyAddOpts.session, "session", "s", "", "session context")
}
