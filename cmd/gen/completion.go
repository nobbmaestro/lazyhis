package gen

import (
	"os"

	"github.com/spf13/cobra"
)

var CompletionCmd = &cobra.Command{
	Use:       "completion [bash|zsh|fish|powershell]",
	Short:     "Generate the autocompletion script for the specified shell",
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE:      runCompletion,
}

func runCompletion(cmd *cobra.Command, args []string) error {
	root := cmd.Root()

	switch args[0] {
	case "bash":
		return root.GenBashCompletion(os.Stdout)
	case "zsh":
		return root.GenZshCompletion(os.Stdout)
	case "fish":
		return root.GenFishCompletion(os.Stdout, true)
	case "powershell":
		return root.GenPowerShellCompletionWithDesc(os.Stdout)
	}
	return nil
}

func init() {}
