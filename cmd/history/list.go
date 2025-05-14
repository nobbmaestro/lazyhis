package history

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/spf13/cobra"
)

var historyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all history records",
	RunE:  runHistoryList,
}

func runHistoryList(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	app := reg.GetApp()

	records, err := app.Service.GetAllHistory()
	if err != nil {
		return err
	}

	for _, record := range records {
		fmt.Println(record.Command.Command)
	}

	return nil
}
