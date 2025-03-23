package history

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/spf13/cobra"
)

var historyLastCmd = &cobra.Command{
	Use:   "last",
	Short: "Last added history record",
	RunE:  runHistoryLast,
}

func runHistoryLast(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)

	record, err := historyService.GetLastHistory()
	if err != nil {
		return err
	}
	fmt.Println(record.Command.Command)

	return nil
}
