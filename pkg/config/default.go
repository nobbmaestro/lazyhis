package config

import path "github.com/nobbmaestro/lazyhis/pkg/config/parsers"

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
			InitialFilterMode: NoFilter,
			CyclicFilterModes: []FilterMode{
				NoFilter,
				PathFilter,
				SessionFilter,
				PathSessionFilter,
				ExitFilter,
			},
			PersistentFilterModes: []FilterMode{
				UniqueFilter,
			},
		},
		Os: OsConfig{
			FetchCurrentSessionCmd: "tmux display-message -p '#S'",
		},
		Log: LogConfig{
			LogEnabled: false,
			LogLevel:   LevelError,
			LogFile:    path.New("~/Library/Logs/lazyhis.log"),
		},
	}
}
