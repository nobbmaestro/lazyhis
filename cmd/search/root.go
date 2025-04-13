package search

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/pkg/formatters"
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
	formatString        string
}

var searchOpts = &SearchOptions{}

const LongDiscription = `Non-interactive history search

Available format placeholders:
  {ID}            Unique record ID
  {COMMAND}       The command string
  {EXECUTED_AT}   Timestamp of command execution (UNIX)
  {EXECUTED_IN}   Duration of command execution
  {EXIT_CODE}     Exit code of the command
  {PATH}          Path context of the command
  {SESSION}       Session context of the command
`

const example = `  lazyhis search --unique -- git add
  lazyhis search -f "{ID}:{EXECUTED_AT};{COMMAND}"
  lazyhis search -u -l 10 -e 0 -p ~/repos/lazyhis -- make
`

var SearchCmd = &cobra.Command{
	Use:     "search [flags] -- [KEYWORDS...]",
	Short:   "Non-interactive history search",
	Long:    LongDiscription,
	Example: example,
	Args:    cobra.ArbitraryArgs,
	RunE:    runSearch,
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

	formatter := formatters.NewFormatter(
		formatters.WithOptions(formatters.DefaultTuiFormatOptions()),
		formatters.WithFormat(searchOpts.formatString),
	)

	for _, record := range formatter.HistoryToFormatString(records) {
		fmt.Println(record)
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
		StringVarP(&searchOpts.formatString, "format", "f", "{COMMAND}", "format of the search results")
}
