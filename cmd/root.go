package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/nobbmaestro/lazyhis/pkg/config"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var (
	printDefaultConfig bool
)

var rootCmd = &cobra.Command{
	Use:   "lazyhis",
	Short: "lazyhis",
	RunE:  runRoot,
}

func runRoot(cmd *cobra.Command, args []string) error {
	switch {
	case printDefaultConfig:
		return runPrintDefaultConfig()
	default:
		return nil
	}
}

func runPrintDefaultConfig() error {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	err := encoder.Encode(config.GetDefaultUserConfig())
	if err != nil {
		return fmt.Errorf("Failed to encode default config: %w", err)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf(
		"%s (Built on %s from Git SHA %s)",
		version,
		date,
		commit,
	)
}

func SetContext(ctx context.Context) {
	rootCmd.SetContext(ctx)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.
		Flags().
		BoolVarP(&printDefaultConfig, "config", "c", false, "print the default config")
}
