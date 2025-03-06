package search

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/pkg/context"
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

var SearchCmd = &cobra.Command{
	Use:   "search [KEYWORDS...]",
	Short: "Interactive history search",
	Args:  cobra.ArbitraryArgs,
	Run:   runSearch,
}

func runSearch(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)

	searchNonInteractive(*historyService, args)
}

func searchNonInteractive(
	historyService service.HistoryService,
	args []string,
) {
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
		if record.Command != nil {
			fmt.Println(record.Command.Command)
		}
	}
}

func init() {
	SearchCmd.
		Flags().
		IntVarP(&searchOpts.exitCode, "exit-code", "e", -1, "filter search results by exit code")
	SearchCmd.
		Flags().
		StringVarP(&searchOpts.tmuxSession, "tmux-session", "s", "", "filter search results by tmux session")
	SearchCmd.
		Flags().
		StringVarP(&searchOpts.path, "path", "p", "", "filter search results by path")
	SearchCmd.
		Flags().
		IntVarP(&searchOpts.maxNumSearchResults, "limit", "l", -1, "limit the number of search results")
	SearchCmd.
		Flags().
		IntVarP(&searchOpts.offsetSearchResults, "offset", "o", -1, "offset of the search results")
}
