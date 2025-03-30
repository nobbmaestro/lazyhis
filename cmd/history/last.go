package history

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/spf13/cobra"
)

var historyLastCmd = &cobra.Command{
	Use:   "last",
	Short: "Last added history record",
	RunE:  runHistoryLast,
}

func runHistoryLast(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	svc := reg.GetService()

	record, err := svc.GetLastHistory()
	if err != nil {
		return err
	}
	fmt.Println(record.Command.Command)

	return nil
}
