package gen

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var dstDir string

var ManCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man pages",
	Args:  cobra.ArbitraryArgs,
	RunE:  runGenDocs,
}

func runGenDocs(cmd *cobra.Command, args []string) error {
	header := &doc.GenManHeader{
		Title:   "LAZYHIS",
		Section: "1",
	}

	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return err
	}

	return doc.GenManTree(cmd.Root(), header, dstDir)
}

func init() {
	ManCmd.
		Flags().
		StringVarP(&dstDir, "dst", "d", "./man", "output directory for man pages")
}
