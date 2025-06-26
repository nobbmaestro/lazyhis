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
			ColumnLabels: map[Column]string{
				ColumnCommand:    "Command",
				ColumnExecutedAt: "Executed",
				ColumnExecutedIn: "Duration",
				ColumnExitCode:   "Exit",
				ColumnID:         "ID",
				ColumnPath:       "Path",
				ColumnSession:    "Session",
			},
			ColumnWidths: map[Column]int{
				ColumnCommand:    100,
				ColumnExecutedAt: 10,
				ColumnExecutedIn: 10,
				ColumnExitCode:   5,
				ColumnID:         5,
				ColumnPath:       25,
				ColumnSession:    25,
			},
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
			Keys: GuiKeys{
				AcceptSelected:  []string{"enter"},
				PrefillSelected: []string{"ctrl+o"},
				DeleteSelected:  []string{"ctrl+x"},
				CopySelected:    []string{"ctrl+y"},
				NextFilter:      []string{"tab"},
				PrevFilter:      []string{"shift+tab"},
				JumpDown:        []string{"ctrl+d"},
				JumpUp:          []string{"ctrl+u"},
				MoveDown:        []string{"ctrl+n", "down"},
				MoveUp:          []string{"ctrl+p", "up"},
				Quit:            []string{"ctrl+q", "ctrl+c", "esc"},
				ShowHelp:        []string{"?"},
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
