package search

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/spf13/cobra"
)

type SearchOptions struct {
	exitCode            int
	maxNumSearchResults int
	offsetSearchResults int
	path                string
	session             string
	uniqueSearchResults bool
}

var searchOpts = &SearchOptions{}

var SearchCmd = &cobra.Command{
	Use:   "search [flags] -- [KEYWORDS...]",
	Short: "Non-interactive history search",
	Args:  cobra.ArbitraryArgs,
	RunE:  runSearch,
}

func runSearch(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	svc := reg.GetService()

	records, err := svc.SearchHistory(
		args,
		searchOpts.exitCode,
		searchOpts.path,
		searchOpts.session,
		searchOpts.maxNumSearchResults,
		searchOpts.offsetSearchResults,
		searchOpts.uniqueSearchResults,
	)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Command != nil {
			fmt.Println(record.Command.Command)
		}
	}

	return nil
}

func init() {
	SearchCmd.
		Flags().
		IntVarP(&searchOpts.exitCode, "exit-code", "e", -1, "filter search results by exit code")
	SearchCmd.
		Flags().
		StringVarP(&searchOpts.session, "session", "s", "", "filter search results by session")
	SearchCmd.
		Flags().
		StringVarP(&searchOpts.path, "path", "p", "", "filter search results by path")
	SearchCmd.
		Flags().
		IntVarP(&searchOpts.maxNumSearchResults, "limit", "l", -1, "limit the number of search results")
	SearchCmd.
		Flags().
		BoolVarP(&searchOpts.uniqueSearchResults, "unique", "u", false, "filter search results by unique commands")
	SearchCmd.
		Flags().
		IntVarP(&searchOpts.offsetSearchResults, "offset", "o", -1, "offset of the search results")
}
