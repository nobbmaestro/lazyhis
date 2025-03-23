package history

import (
	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/spf13/cobra"
)

var historyPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Delete history records matching the configured exclusion filters",
	RunE:  runHistoryPrune,
}

func runHistoryPrune(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)

	err := historyService.PruneHistory()
	if err != nil {
		return err
	}

	return nil
}
