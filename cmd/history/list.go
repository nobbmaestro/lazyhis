package history

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/spf13/cobra"
)

var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all history records",
	RunE:  runHistoryList,
}

func runHistoryList(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)

	records, err := historyService.GetAllHistory()
	if err != nil {
		return err
	}

	for _, record := range records {
		fmt.Println(record.Command.Command)
	}

	return nil
}
