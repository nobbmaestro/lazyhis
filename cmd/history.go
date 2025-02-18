package cmd

import (
	"fmt"
	"os"

	"github.com/nobbmaestro/lazyhis/pkg/context"
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

var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all items in history",
	Run:   runHistoryList,
}

var historyPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Delete history entries matching the configured exclusion filters",
	Run:   runHistoryPrune,
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
	historyCmd.AddCommand(historyListCmd)
	historyCmd.AddCommand(historyPruneCmd)
	rootCmd.AddCommand(historyCmd)
}

func runHistoryAdd(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)
	config := context.GetConfig(ctx)

	_, err := historyService.AddHistory(
		args,
		&historyOpts.exitCode,
		&historyOpts.executedIn,
		&historyOpts.path,
		&historyOpts.tmuxSession,
		&config.Db.ExcludeCommands,
	)
	if err != nil {
		return
	}
}

func runHistoryList(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)

	commands, err := historyService.GetAllCommands()
	if err != nil {
		return
	}

	for _, command := range commands {
		fmt.Println(command.Command)
	}
}

func runHistoryPrune(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)
	config := context.GetConfig(ctx)

	err := historyService.PruneHistory(config.Db.ExcludeCommands)
	if err != nil {
		return
	}
}
