package gen

import "github.com/spf13/cobra"

var GenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate auxiliary files such as docs and shell completion",
}

func init() {
	GenCmd.AddCommand(
		CompletionCmd,
		ManCmd,
	)
}
