package history

import (
	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/spf13/cobra"
)

var historyPruneFlags = &HistoryFlags{}

var historyPruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Delete history records matching the configured exclusion filters",
	RunE:  runHistoryPrune,
}

func runHistoryPrune(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	svc := reg.GetService()

	err := svc.PruneHistory(
		historyPruneFlags.dryRun,
		historyPruneFlags.verbose,
	)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	historyPruneCmd.
		Flags().
		BoolVarP(&historyPruneFlags.dryRun, "dry", "d", false, "list matching commands without performing actual deletion")
	historyPruneCmd.
		Flags().
		BoolVarP(&historyPruneFlags.verbose, "verbose", "v", false, "increase verbosity")
}
