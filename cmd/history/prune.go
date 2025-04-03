package history

import (
	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/spf13/cobra"
)

var historyPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Delete history records matching the configured exclusion filters",
	RunE:  runHistoryPrune,
}

func runHistoryPrune(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	svc := reg.GetService()

	err := svc.PruneHistory()
	if err != nil {
		return err
	}

	return nil
}
