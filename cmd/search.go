package cmd

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/domain/service"
	"github.com/spf13/cobra"
)

type SearchOptions struct {
	exitCode            int
	maxNumSearchResults int
	offsetSearchResults int
	path                string
	tmuxSession         string
}

var searchOpts = &SearchOptions{}

var searchCmd = &cobra.Command{
	Use:   "search [KEYWORDS...]",
	Short: "Interactive history search",
	Args:  cobra.ArbitraryArgs,
	Run:   runSearch,
}

func runSearch(cmd *cobra.Command, args []string) {
	serviceCtx := cmd.Context().Value(service.ServiceCtxKey).(*service.ServiceContext)
	historyService := serviceCtx.HistoryService

	records, err := historyService.SearchHistory(
		args,
		searchOpts.exitCode,
		searchOpts.path,
		searchOpts.tmuxSession,
		searchOpts.maxNumSearchResults,
		searchOpts.offsetSearchResults,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, record := range records {
		fmt.Println(record.Command.Command)
	}
}

func init() {
	searchCmd.
		Flags().
		IntVarP(&searchOpts.exitCode, "exit-code", "e", -1, "Filter search results by exit code")
	searchCmd.
		Flags().
		StringVarP(&searchOpts.tmuxSession, "tmux-session", "s", "", "Filter search results by tmux session")
	searchCmd.
		Flags().
		StringVarP(&searchOpts.path, "path", "p", "", "Filter search results by path")
	searchCmd.
		Flags().
		IntVarP(&searchOpts.maxNumSearchResults, "limit", "l", -1, "Limit the number of search results")
	searchCmd.
		Flags().
		IntVarP(&searchOpts.offsetSearchResults, "offset", "o", -1, "Offset of the search results")

	rootCmd.AddCommand(searchCmd)
}
