package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nobbmaestro/lazyhis/cmd/history"
	"github.com/nobbmaestro/lazyhis/cmd/initialize"
	"github.com/nobbmaestro/lazyhis/cmd/search"
	"github.com/nobbmaestro/lazyhis/pkg/config"
	"github.com/nobbmaestro/lazyhis/pkg/domain/model"
	"github.com/nobbmaestro/lazyhis/pkg/gui"
	"github.com/nobbmaestro/lazyhis/pkg/gui/formatters"
	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/nobbmaestro/lazyhis/pkg/utils"
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
		return runHistoryGui(cmd, args)
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

func runHistoryGui(
	cmd *cobra.Command,
	args []string,
) error {
	reg := registry.FromContext(cmd.Context())
	cfg := reg.GetConfig()
	svc := reg.GetService()

	partialSearchHistory := func(keywords []string, mode config.FilterMode) []model.History {
		records, err := svc.SearchHistory(
			keywords,
			applyExitCodeFilter(mode, cfg.Gui.PersistentFilterModes),
			applyPathFilter(mode, cfg.Gui.PersistentFilterModes),
			applySessionFilter(
				mode,
				cfg.Gui.PersistentFilterModes,
				cfg.Os.FetchCurrentSessionCmd,
			),
			-1, //maxNumSearchResults
			-1, //offsetSearchResults
			applyUniqueCommandFilter(mode, cfg.Gui.PersistentFilterModes),
		)
		if err != nil {
			return nil
		}
		return records
	}

	p := tea.NewProgram(
		gui.NewGui(
			partialSearchHistory,
			cfg.Gui,
			gui.WithVersion(cmd.Version),
			gui.WithInitialQuery(args),
			gui.WithFormatter(
				formatters.NewFmt(formatters.WithColumns(cfg.Gui.ColumnLayout)),
			),
		),
		tea.WithAltScreen(),
	)

	result, err := p.Run()
	if err != nil {
		return err
	}

	if model, ok := result.(*gui.Model); ok && model.SelectedRecord.Command != nil {
		command := model.SelectedRecord.Command.Command
		switch model.UserAction {
		case gui.ActionAcceptSelected:
			fmt.Fprintf(os.Stderr, "__lazyhis_accept__:%s\n", command)
		case gui.ActionPrefillSelected:
			fmt.Fprintf(os.Stderr, "__lazyhis_prefill__:%s\n", command)
		}
	}

	return nil
}

func applyPathFilter(
	mode config.FilterMode,
	persistent []config.FilterMode,
) string {
	if mode == config.WorkdirFilter ||
		mode == config.WorkdirSessionFilter ||
		slices.Contains(persistent, config.WorkdirFilter) ||
		slices.Contains(persistent, config.WorkdirSessionFilter) {
		if p, err := os.Getwd(); err == nil {
			return p
		}
	}
	return ""
}

func applySessionFilter(
	mode config.FilterMode,
	persistent []config.FilterMode,
	sessionCmd string,
) string {
	if mode == config.SessionFilter ||
		mode == config.WorkdirSessionFilter ||
		slices.Contains(persistent, config.SessionFilter) {
		if s, err := utils.RunCommand(strings.Split(sessionCmd, " ")); err == nil {
			return s
		}
	}
	return ""
}

func applyExitCodeFilter(
	mode config.FilterMode,
	persistent []config.FilterMode,
) int {
	if mode == config.SuccessFilter ||
		slices.Contains(persistent, config.SuccessFilter) {
		return 0
	}
	return -1
}

func applyUniqueCommandFilter(
	mode config.FilterMode,
	persistent []config.FilterMode,
) bool {
	if mode == config.UniqueFilter || slices.Contains(persistent, config.UniqueFilter) {
		return true
	}
	return false
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
