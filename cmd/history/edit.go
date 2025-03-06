package history

import (
	"strconv"

	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/spf13/cobra"
)

var historyEditOpts = &HistoryOptions{}

var historyEditCmd = &cobra.Command{
	Use:   "edit [ID]",
	Short: "Edit history record by ID",
	Args:  cobra.ExactArgs(1),
	Run:   runHistoryEdit,
}

func runHistoryEdit(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)

	historyID, err := strconv.Atoi(args[0])
	if err != nil {
		return
	}

	var (
		exitCode    *int
		executedIn  *int
		path        *string
		tmuxSession *string
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
	if cmd.Flags().Changed("tmux-session") {
		tmuxSession = &historyEditOpts.tmuxSession
	}

	_, err = historyService.EditHistory(
		historyID,
		exitCode,
		executedIn,
		path,
		tmuxSession,
	)
	if err != nil {
		return
	}
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
		StringVarP(&historyEditOpts.tmuxSession, "tmux-session", "s", "", "tmux session context")
}
