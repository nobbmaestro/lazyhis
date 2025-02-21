package cmd

import (
	"fmt"
	"os"
	"strconv"

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

var historyAddOpts = &HistoryOptions{}
var historyEditOpts = &HistoryOptions{}

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

var historyEditCmd = &cobra.Command{
	Use:   "edit [ID]",
	Short: "Edit history by ID",
	Args:  cobra.ExactArgs(1),
	Run:   runHistoryEdit,
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

	historyCmd.AddCommand(
		historyAddCmd,
		historyEditCmd,
		historyListCmd,
		historyPruneCmd,
	)
	rootCmd.AddCommand(historyCmd)
}
