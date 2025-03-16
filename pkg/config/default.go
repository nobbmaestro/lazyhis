package config

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
	}
}
