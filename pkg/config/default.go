package config

import path "github.com/nobbmaestro/lazyhis/pkg/config/parsers"

func GetDefaultUserConfig() *UserConfig {
	return &UserConfig{
		Db: DbConfig{
			ExcludeCommands: []string{},
			ExcludePrefix:   "\x20",
		},
		Gui: GuiConfig{
			ShowColumnLabels: false,
			ColumnLayout: []Column{
				ColumnExitCode,
				ColumnExecutedAt,
				ColumnCommand,
			},
			InitialFilterMode: NoFilter,
			CyclicFilterModes: []FilterMode{
				NoFilter,
				WorkdirFilter,
				SessionFilter,
				WorkdirSessionFilter,
				SuccessFilter,
			},
			PersistentFilterModes: []FilterMode{
				UniqueFilter,
			},
			Theme: GuiTheme{
				BorderColor:        "8",
				FilterFgColor:      "7",
				HelpAccentColor:    "6",
				HelpFgColor:        "7",
				InputFgColor:       "7",
				TableCursorBgColor: "",
				TableCursorFgColor: "6",
				TableLabelsFgColor: "7",
				VersionFgColor:     "6",
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
