package history

import (
	"github.com/spf13/cobra"
)

type HistoryOptions struct {
	executedIn int
	exitCode   int
	path       string
	session    string
}

var HistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Manipulate history database",
}

func init() {
	HistoryCmd.AddCommand(
		historyAddCmd,
		historyEditCmd,
		historyImportCmd,
		historyLastCmd,
		historyListCmd,
		historyPruneCmd,
	)
}
