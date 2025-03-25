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
	"github.com/nobbmaestro/lazyhis/pkg/utils"
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
	cfg config.UserConfig,
	args []string,
	version string,
) error {
	partialSearchHistory := func(keywords []string, mode config.FilterMode) []model.History {
		records, err := historyService.SearchHistory(
			append(args, keywords...),
			applyExitCodeFilter(mode),
			applyPathFilter(mode),
			applySessionFilter(mode, cfg.Os.FetchCurrentSessionCmd),
			-1, //maxNumSearchResults
			-1, //offsetSearchResults
			cfg.Gui.ShowUniqueCommands,
		)
		if err != nil {
			return nil
		}
		return records
	}

	p := tea.NewProgram(
		gui.NewModel(
			cfg.Gui,
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

func applyPathFilter(mode config.FilterMode) string {
	if mode == config.PathFilter || mode == config.PathSessionFilter {
		if p, err := os.Getwd(); err == nil {
			return p
		}
	}
	return ""
}

func applySessionFilter(mode config.FilterMode, sessionCmd string) string {
	if mode == config.SessionFilter || mode == config.PathSessionFilter {
		if s, err := utils.RunCommand(strings.Split(sessionCmd, " ")); err == nil {
			return s
		}
	}
	return ""
}

func applyExitCodeFilter(mode config.FilterMode) int {
	if mode == config.ExitFilter {
		return 0
	}
	return -1
}

func init() {
	SearchCmd.
		Flags().
		IntVarP(&searchOpts.exitCode, "exit-code", "e", -1, "filter search results by exit code (non-interactive only)")
	SearchCmd.
		Flags().
		StringVarP(&searchOpts.session, "session", "s", "", "filter search results by session (non-interactive only)")
	SearchCmd.
		Flags().
		StringVarP(&searchOpts.path, "path", "p", "", "filter search results by path (non-interactive only)")
	SearchCmd.
		Flags().
		IntVarP(&searchOpts.maxNumSearchResults, "limit", "l", -1, "limit the number of search results (non-interactive only)")
	SearchCmd.
		Flags().
		BoolVarP(&searchOpts.uniqueSearchResults, "unique", "u", false, "filter search results by unique commands (non-interactive only)")
	SearchCmd.
		Flags().
		IntVarP(&searchOpts.offsetSearchResults, "offset", "o", -1, "offset of the search results (non-interactive only)")
	SearchCmd.
		Flags().
		BoolVarP(&searchOpts.runInteractive, "interactive", "i", false, "open interactive search GUI")
}
