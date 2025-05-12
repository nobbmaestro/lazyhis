package history

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	appopts "github.com/nobbmaestro/lazyhis/pkg/app"

	"github.com/nobbmaestro/lazyhis/pkg/registry"
	"github.com/spf13/cobra"
)

var isZshHistfile bool
var historyImportFlags = &HistoryFlags{}

var historyImportCmd = &cobra.Command{
	Use:   "import [HISTFILE]",
	Short: "Import history from histfile",
	Args:  cobra.ExactArgs(1),
	RunE:  runImport,
}

func runImport(cmd *cobra.Command, args []string) error {
	switch {
	case isZshHistfile:
		return importZshHistfile(cmd, args)
	default:
		return fmt.Errorf("Unsupported shell option")
	}
}

func importZshHistfile(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry(registry.WithContext(cmd.Context()))
	app := reg.GetApp()

	file, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := parseHistoryLine(scanner.Text())
		if len(command) == 0 {
			continue
		}

		_, err = app.AddHistory(
			historyImportFlags.dryRun,
			historyImportFlags.verbose,
			true, // AddUniqueOnly
			appopts.WithQuery(command),
		)
		if err != nil {
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error reading file: %w", err)
	}

	return nil
}

func parseHistoryLine(line string) []string {
	if strings.HasPrefix(line, ": ") {
		parts := strings.SplitN(line, ";", 2)
		if len(parts) < 2 {
			return []string{}
		}
		line = parts[1]
	}

	rePatterns := []*regexp.Regexp{
		regexp.MustCompile(`[\s&&[^\x20]]`), // whitespaces, except space
		regexp.MustCompile(`[^\x20-\x7E]`),  // non-printable ASCII
	}
	for _, re := range rePatterns {
		line = re.ReplaceAllString(line, "")
	}

	return strings.Fields(line)
}

func init() {
	historyImportCmd.
		Flags().
		BoolVar(&isZshHistfile, "zsh", false, "import zsh histfile")
	historyImportCmd.
		Flags().
		BoolVarP(&historyImportFlags.dryRun, "dry", "d", false, "list commands without performing actual import")
	historyImportCmd.
		Flags().
		BoolVarP(&historyImportFlags.verbose, "verbose", "v", false, "increase verbosity")
}
