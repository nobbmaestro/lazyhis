package initialize

import (
	"fmt"

	"github.com/nobbmaestro/lazyhis/scripts"
	"github.com/spf13/cobra"
)

var initZshCmd = &cobra.Command{
	Use: "zsh",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(scripts.InitZsh)
	},
}
