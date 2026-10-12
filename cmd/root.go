package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nobbmaestro/lazyhis/cmd/gen"
	"github.com/nobbmaestro/lazyhis/cmd/history"
	"github.com/nobbmaestro/lazyhis/cmd/initialize"
	"github.com/nobbmaestro/lazyhis/cmd/search"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/formatters"
	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/nobbmaestro/lazyhis/pkg/tui"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var (
	printUserConfig    bool
	printDefaultConfig bool
	printConfigPath    bool
)

var rootCmd = &cobra.Command{
	Use:   "lazyhis [flags] -- [KEYWORDS...]",
	Short: "lazyhis",
	Args:  cobra.ArbitraryArgs,
	RunE:  runRoot,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func runRoot(cmd *cobra.Command, args []string) error {
	switch {
	case printUserConfig:
		return runPrintUserConfig(cmd, args)
	case printDefaultConfig:
		return runPrintDefaultConfig(cmd, args)
	case printConfigPath:
		return runPrintConfigPath(cmd, args)
	default:
		return runHistoryTui(cmd, args)
	}
}

func runPrintUserConfig(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	return printConfig(reg.GetConfig())
}

func runPrintDefaultConfig(cmd *cobra.Command, args []string) error {
	return printConfig(config.GetDefaultUserConfig())
}

func printConfig(cfg *config.UserConfig) error {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("error encoding config: %w", err)
	}
	fmt.Printf("%s\n", buf.String())
	return nil
}

func runPrintConfigPath(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	fmt.Println(reg.GetConfigPath())
	return nil
}

func runHistoryTui(
	cmd *cobra.Command,
	args []string,
) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	cfg := reg.GetConfig()
	app := reg.GetApp()

	p := tea.NewProgram(
		tui.NewTui(
			app,
			&cfg.Ui,
			tui.WithInitialQuery(args),
			tui.WithFormatter(
				formatters.NewFormatter(
					formatters.WithColumns(cfg.Ui.ColumnLayout),
					formatters.WithOptions(formatters.DefaultUiFormatOptions()),
				),
			),
		),
		tea.WithAltScreen(),
	)

	result, err := p.Run()
	if err != nil {
		return err
	}

	if model, ok := result.(*tui.Model); ok && model.SelectedRecord.Command != nil {
		command := model.SelectedRecord.Command.Command
		switch model.ExitCode {
		case tui.ExitAcceptSelected:
			fmt.Fprintf(os.Stderr, "__lazyhis_accept__:%s\n", command)
		case tui.ExitPrefillSelected:
			fmt.Fprintf(os.Stderr, "__lazyhis_prefill__:%s\n", command)
		}
	}

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
		BoolVarP(&printConfigPath, "config-path", "p", false, "print the config file path")

	rootCmd.AddCommand(history.HistoryCmd)
	rootCmd.AddCommand(initialize.InitCmd)
	rootCmd.AddCommand(search.SearchCmd)
	rootCmd.AddCommand(gen.GenCmd)
}
