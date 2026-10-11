package config

import (
	"os"
	"path/filepath"

	path "github.com/nobbmaestro/lazyhis/pkg/config/parsers"
)

var DefaultPath = filepath.Join(
	os.Getenv("HOME"),
	".config",
	"lazyhis",
	"lazyhis.yml")

func ResolveConfigPath() string {
	if p := os.Getenv("LAZYHIS_CONFIG"); p != "" {
		return string(path.New(p))
	}

	if p := os.Getenv("XDG_CONFIG_HOME"); p != "" {
		return filepath.Join(p, "lazyhis", "lazyhis.yml")
	}

	return DefaultPath
}
