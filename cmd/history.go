package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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

var historyLastCmd = &cobra.Command{
	Use:   "last",
	Short: "Last added history record",
	Run:   runHistoryLast,
}

var historyImportCmd = &cobra.Command{
	Use:   "import [HISTFILE]",
	Short: "Import history from histfile",
	Args:  cobra.ExactArgs(1),
	Run:   runImport,
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

func runHistoryLast(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)

	record, err := historyService.GetLastHistory()
	if err != nil {
		return
	}
	fmt.Println(record.Command.Command)
}

func runHistoryList(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)

	records, err := historyService.GetAllHistory()
	if err != nil {
		return
	}

	for _, record := range records {
		fmt.Println(record.Command.Command)
	}
}

func runImport(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)
	config := context.GetConfig(ctx)

	file, err := os.Open(args[0])
	if err != nil {
		fmt.Println(fmt.Errorf("%w", err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := parseHistoryLine(scanner.Text())
		if len(command) == 0 {
			continue
		}

		_, err = historyService.AddHistoryIfUnique(
			command,
			nil,
			nil,
			nil,
			nil,
			&config.Db.ExcludeCommands,
		)
		if err != nil {
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(fmt.Errorf("Error reading file: %w", err))
	}
}

func parseHistoryLine(line string) []string {
	if strings.HasPrefix(line, ": ") {
		parts := strings.SplitN(line, ";", 2)
		if len(parts) < 2 {
			return []string{}
		}
		line = parts[1]
	}

	rePatterns := []*regexp.Regexp{
		regexp.MustCompile(`[\s&&[^\x20]]`), // whitespaces, except space
		regexp.MustCompile(`[^\x20-\x7E]`),  // non-printable ASCII
	}
	for _, re := range rePatterns {
		line = re.ReplaceAllString(line, "")
	}

	return strings.Fields(line)
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
		historyImportCmd,
		historyLastCmd,
		historyListCmd,
		historyPruneCmd,
	)
	rootCmd.AddCommand(historyCmd)
}
