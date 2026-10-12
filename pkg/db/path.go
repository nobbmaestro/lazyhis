package db

import (
	"os"
	"path/filepath"

	path "github.com/nobbmaestro/lazyhis/pkg/config/parsers"
)

var DefaultPath = filepath.Join(
	os.Getenv("HOME"),
	".lazyhis.db")

func ResolveDbPath() string {
	if p := os.Getenv("LAZYHIS_DB"); p != "" {
		return string(path.New(p))
	}

	if p := os.Getenv("XDG_DATA_HOME"); p != "" {
		return filepath.Join(p, "lazyhis", "lazyhis.db")
	}

	return DefaultPath
}
