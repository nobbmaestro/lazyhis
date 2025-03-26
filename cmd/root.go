package cmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/nobbmaestro/lazyhis/cmd/history"
	"github.com/nobbmaestro/lazyhis/cmd/initialize"
	"github.com/nobbmaestro/lazyhis/cmd/search"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var (
	printUserConfig    bool
	printDefaultConfig bool
	printConfigPath    bool
)

var rootCmd = &cobra.Command{
	Use:   "lazyhis",
	Short: "lazyhis",
	RunE:  runRoot,
}

func runRoot(cmd *cobra.Command, args []string) error {
	switch {
	case printUserConfig:
		return runPrintUserConfig()
	case printDefaultConfig:
		return runPrintDefaultConfig()
	case printConfigPath:
		return runPrintConfigPath()
	default:
		return nil
	}
}

func runPrintUserConfig() error {
	cfg, err := config.ReadUserConfig()
	if err != nil {
		return fmt.Errorf("Failed to read user config: %w", err)
	}
	return printConfig(cfg)
}

func runPrintDefaultConfig() error {
	return printConfig(config.GetDefaultUserConfig())
}

func printConfig(cfg *config.UserConfig) error {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("Error encoding config: %w", err)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

func runPrintConfigPath() error {
	fmt.Println(config.GetUserConfigPath())
	return nil
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = version
}

func SetContext(ctx context.Context) {
	rootCmd.SetContext(ctx)
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.
		Flags().
		BoolVarP(&printUserConfig, "user-config", "c", false, "print the user config")
	rootCmd.
		Flags().
		BoolVarP(&printDefaultConfig, "default-config", "C", false, "print the default config")
	rootCmd.
		Flags().
		BoolVarP(&printConfigPath, "config-dir", "d", false, "print the config directory")

	rootCmd.AddCommand(history.HistoryCmd)
	rootCmd.AddCommand(initialize.InitCmd)
	rootCmd.AddCommand(search.SearchCmd)
}
