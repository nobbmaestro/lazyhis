package history

import (
	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/spf13/cobra"
)

var historyPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Delete history records matching the configured exclusion filters",
	Run:   runHistoryPrune,
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
