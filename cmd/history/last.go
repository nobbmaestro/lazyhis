package history

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/spf13/cobra"
)

var historyLastCmd = &cobra.Command{
	Use:   "last",
	Short: "Last added history record",
	Run:   runHistoryLast,
}

func runHistoryLast(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)

	record, err := historyService.GetLastHistory()
	if err != nil {
		return
	}
	fmt.Println(record.Command.Command)
}
