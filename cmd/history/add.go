package history

import (
	"fmt"
	"os"

	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
	"github.com/spf13/cobra"
)

var historyAddOpts = &HistoryOptions{}

var historyAddCmd = &cobra.Command{
	Use:   "add [CMD...]",
	Short: "Add history record",
	Args:  cobra.ArbitraryArgs,
	Run:   runHistoryAdd,
}

func runHistoryAdd(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)
	config := context.GetConfig(ctx)

	record, err := historyService.AddHistory(
		args,
		&historyAddOpts.exitCode,
		&historyAddOpts.executedIn,
		&historyAddOpts.path,
		&historyAddOpts.tmuxSession,
		&config.Db.ExcludeCommands,
	)
	if err != nil {
		return
	}
	if record != nil {
		fmt.Println(record.ID)
	}
}

func init() {
	currentPath, err := os.Getwd()
	if err != nil {
		currentPath = ""
	}
	currentTmuxSession, err := utils.GetCurrentTmuxSession()
	if err != nil {
		currentTmuxSession = ""
	}

	historyAddCmd.
		Flags().
		IntVarP(&historyAddOpts.exitCode, "exit-code", "e", 0, "exit code for the command")
	historyAddCmd.
		Flags().
		IntVarP(&historyAddOpts.executedIn, "duration", "d", 0, "execution duration of the CMD in milliseconds")
	historyAddCmd.
		Flags().
		StringVarP(&historyAddOpts.path, "path", "p", currentPath, "working directory context")
	historyAddCmd.
		Flags().
		StringVarP(&historyAddOpts.tmuxSession, "tmux-session", "s", currentTmuxSession, "tmux session context")
}
