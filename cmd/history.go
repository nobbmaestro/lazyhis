package cmd

import (
	"os"

	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
	"github.com/spf13/cobra"
)

type HistoryOptions struct {
	executedIn  int
	exitCode    int
	path        string
	tmuxSession string
}

var historyOpts = &HistoryOptions{}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Manipulate history database",
}

var historyAddCmd = &cobra.Command{
	Use:   "add [CMD...]",
	Short: "Add to history",
	Args:  cobra.ArbitraryArgs,
	Run:   runHistoryAdd,
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
		IntVarP(&historyOpts.exitCode, "exit-code", "e", 0, "Exit code for the command")
	historyAddCmd.
		Flags().
		StringVarP(&historyOpts.path, "path", "p", currentPath, "Working directory context")
	historyAddCmd.
		Flags().
		StringVarP(&historyOpts.tmuxSession, "tmux-session", "s", currentTmuxSession, "Tmux session context")

	historyCmd.AddCommand(historyAddCmd)
	rootCmd.AddCommand(historyCmd)
}

func runHistoryAdd(cmd *cobra.Command, args []string) {
	serviceCtx := cmd.Context().Value(service.ServiceCtxKey).(*service.ServiceContext)
	historyService := serviceCtx.HistoryService

	_, err := historyService.AddHistory(
		args,
		&historyOpts.exitCode,
		&historyOpts.executedIn,
		&historyOpts.path,
		&historyOpts.tmuxSession,
	)
	if err != nil {
		return
	}
}
