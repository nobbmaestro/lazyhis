package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nobbmaestro/lazyhis/pkg/context"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [HISTFILE]",
	Short: "Import history from histfile",
	Args:  cobra.ExactArgs(1),
	Run:   runImport,
}

func runImport(cmd *cobra.Command, args []string) {
	ctx := cmd.Context()
	historyService := context.GetService(ctx)
	config := context.GetConfig(ctx)

	file, err := os.Open(args[0])
	if err != nil {
		fmt.Println(fmt.Errorf("%w", err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		command := parseHistoryLine(scanner.Text())
		if len(command) == 0 {
			continue
		}

		_, err = historyService.AddHistoryIfUnique(
			command,
			nil,
			nil,
			nil,
			nil,
			&config.Db.ExcludeCommands,
		)
		if err != nil {
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(fmt.Errorf("Error reading file: %w", err))
	}
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
	rootCmd.AddCommand(importCmd)
}
