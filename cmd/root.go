package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lazyhis",
	Short: "Command history with a database twist",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello, World!")
	},
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
