package history

import (
	"fmt"
	"os"
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
	"github.com/spf13/cobra"
)

var historyAddOpts = &HistoryOptions{}

var historyAddCmd = &cobra.Command{
	Use:   "add [CMD...]",
	Short: "Add history record",
	Args:  cobra.ArbitraryArgs,
	RunE:  runHistoryAdd,
}

func runHistoryAdd(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	cfg := reg.GetConfig()
	svc := reg.GetService()

	if historyAddOpts.path == "" {
		if currentPath, err := os.Getwd(); err == nil {
			historyAddOpts.path = currentPath
		}
	}

	if historyAddOpts.session == "" {
		cmd := strings.Fields(cfg.Os.FetchCurrentSessionCmd)
		if currentSession, err := utils.RunCommand(cmd); err == nil {
			historyAddOpts.session = currentSession
		}
	}

	record, err := svc.AddHistory(
		args,
		&historyAddOpts.exitCode,
		&historyAddOpts.executedIn,
		&historyAddOpts.path,
		&historyAddOpts.session,
		false, // dryRun
		false, // verbose
	)
	if err != nil {
		return err
	}
	if record != nil {
		fmt.Println(record.ID)
	}

	return nil
}

func init() {
	historyAddCmd.
		Flags().
		IntVarP(&historyAddOpts.exitCode, "exit-code", "e", 0, "exit code for the command")
	historyAddCmd.
		Flags().
		IntVarP(&historyAddOpts.executedIn, "duration", "d", 0, "execution duration of the CMD in milliseconds")
	historyAddCmd.
		Flags().
		StringVarP(&historyAddOpts.path, "path", "p", "", "working directory context")
	historyAddCmd.
		Flags().
		StringVarP(&historyAddOpts.session, "session", "s", "", "session context")
}
