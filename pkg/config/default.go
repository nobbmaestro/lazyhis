package config

import (
	"os"
	"path/filepath"
)

func GetDefaultUserConfig() *UserConfig {
	return &UserConfig{
		Db: DbConfig{
			ExcludeCommands: []string{},
			ExcludePrefix:   "\x20",
		},
		Gui: GuiConfig{
			ColumnLayout: []Column{
				ColumnExitCode,
				ColumnExecutedAt,
				ColumnCommand,
			},
			ShowUniqueCommands: true,
		},
		Os: OsConfig{
			FetchCurrentSessionCmd: "tmux display-message -p '#S'",
		},
		Log: LogConfig{
			LogEnabled: false,
			LogLevel:   LevelError,
			LogFile: filepath.Join(
				os.Getenv("HOME"),
				"Library",
				"Logs",
				"lazyhis.log",
			),
		},
	}
}
