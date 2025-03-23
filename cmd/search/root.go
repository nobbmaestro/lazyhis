package search

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/domain/service"
	"github.com/nobbmaestro/lazyhis/pkg/gui"
	"github.com/spf13/cobra"
)

type SearchOptions struct {
	exitCode            int
	maxNumSearchResults int
	offsetSearchResults int
	path                string
	session             string
	runInteractive      bool
	uniqueSearchResults bool
}

var searchOpts = &SearchOptions{}

var SearchCmd = &cobra.Command{
	Use:   "search [KEYWORDS...]",
	Short: "Interactive history search",
	Args:  cobra.ArbitraryArgs,
	RunE:  runSearch,
}

func runSearch(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)
	config := context.GetConfig(ctx)

	if searchOpts.runInteractive {
		return searchInteractive(*historyService, *config, args, cmd.Root().Version)
	}
	return searchNonInteractive(*historyService, args)
}

func searchNonInteractive(
	historyService service.HistoryService,
	args []string,
) error {
	records, err := historyService.SearchHistory(
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

func searchInteractive(
	historyService service.HistoryService,
	config config.UserConfig,
	args []string,
	version string,
) error {
	partialSearchHistory := func(keywords []string) []model.History {
		records, err := historyService.SearchHistory(
			append(args, keywords...),
			searchOpts.exitCode,
			searchOpts.path,
			searchOpts.session,
			searchOpts.maxNumSearchResults,
			searchOpts.offsetSearchResults,
			searchOpts.uniqueSearchResults || config.Gui.ShowUniqueCommands,
		)
		if err != nil {
			return nil
		}
		return records
	}

	p := tea.NewProgram(
		gui.NewModel(
			config.Gui.ColumnLayout,
			partialSearchHistory,
			version,
			strings.Join(args, " "),
		),
		tea.WithAltScreen(),
	)

	result, err := p.Run()
	if err != nil {
		return err
	}

	if model, ok := result.(*gui.Model); ok && model.SelectedRecord.Command != nil {
		command := model.SelectedRecord.Command.Command
		switch model.UserAction {
		case gui.ActionAcceptSelected:
			fmt.Fprintf(os.Stderr, "__lazyhis_accept__:%s\n", command)
		case gui.ActionPrefillSelected:
			fmt.Fprintf(os.Stderr, "__lazyhis_prefill__:%s\n", command)
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
	SearchCmd.
		Flags().
		BoolVarP(&searchOpts.runInteractive, "interactive", "i", false, "open interactive search GUI")
}
