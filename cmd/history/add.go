package history

import (
	"fmt"
	"os"
	"strings"

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

	if historyAddOpts.path == "" {
		if currentPath, err := os.Getwd(); err == nil {
			historyAddOpts.path = currentPath
		}
	}

	if historyAddOpts.session == "" {
		cmd := strings.Fields(config.Os.FetchCurrentSessionCmd)
		if currentSession, err := utils.RunCommand(cmd); err == nil {
			historyAddOpts.session = currentSession
		}
	}

	record, err := historyService.AddHistory(
		args,
		&historyAddOpts.exitCode,
		&historyAddOpts.executedIn,
		&historyAddOpts.path,
		&historyAddOpts.session,
	)
	if err != nil {
		return
	}
	if record != nil {
		fmt.Println(record.ID)
	}
}

func init() {
	historyAddCmd.
		Flags().
		IntVarP(&historyAddOpts.exitCode, "exit-code", "e", 0, "exit code for the command")
	historyAddCmd.
		Flags().
		IntVarP(&historyAddOpts.executedIn, "duration", "d", 0, "execution duration of the CMD in milliseconds")
	historyAddCmd.
		Flags().
		StringVarP(&historyAddOpts.path, "path", "p", "", "working directory context")
	historyAddCmd.
		Flags().
		StringVarP(&historyAddOpts.session, "session", "s", "", "session context")
}
