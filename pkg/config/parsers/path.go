package path

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Path string

func New(path string) Path {
	return Path(expandPath(path))
}

func (p *Path) UnmarshalYAML(value *yaml.Node) error {
	*p = New(value.Value)
	return nil
}

func (p Path) String() string {
	return string(p)
}

func expandPath(path string) string {
	expandedPath := os.ExpandEnv(path)

	if strings.HasPrefix(expandedPath, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			expandedPath = filepath.Join(home, strings.TrimPrefix(expandedPath, "~"))
		}
	}

	return filepath.Clean(expandedPath)
}
