package initialize

import (
	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Print shell init script",
}

func init() {
	InitCmd.AddCommand(initZshCmd)
}
