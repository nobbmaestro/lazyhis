package history

import (
	"strconv"

	appopts "github.com/nobbmaestro/lazyhis/pkg/app"

	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/spf13/cobra"
)

var historyEditOpts = &HistoryOptions{}

var historyEditCmd = &cobra.Command{
	Use:   "edit [ID]",
	Short: "Edit history record by ID",
	Args:  cobra.ExactArgs(1),
	RunE:  runHistoryEdit,
}

func runHistoryEdit(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	app := reg.GetApp()

	historyID, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	var (
		exitCode   *int
		executedIn *int
		path       *string
		session    *string
	)

	if cmd.Flags().Changed("exit-code") {
		exitCode = &historyEditOpts.exitCode
	}
	if cmd.Flags().Changed("duration") {
		executedIn = &historyEditOpts.executedIn
	}
	if cmd.Flags().Changed("path") {
		path = &historyEditOpts.path
	}
	if cmd.Flags().Changed("session") {
		session = &historyEditOpts.session
	}

	_, err = app.EditHistory(
		historyID,
		appopts.WithExitCode(*exitCode),
		appopts.WithExecutedIn(*executedIn),
		appopts.WithPath(*path),
		appopts.WithSession(*session),
	)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	historyEditCmd.
		Flags().
		IntVarP(&historyEditOpts.exitCode, "exit-code", "e", -1, "exit code for the command")
	historyEditCmd.
		Flags().
		IntVarP(&historyEditOpts.executedIn, "duration", "d", -1, "execution duration of the CMD in microseconds")
	historyEditCmd.
		Flags().
		StringVarP(&historyEditOpts.path, "path", "p", "", "working directory context")
	historyEditCmd.
		Flags().
		StringVarP(&historyEditOpts.session, "session", "s", "", "session context")
}
