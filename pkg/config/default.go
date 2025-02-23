package config

func GetDefaultUserConfig() *UserConfig {
	return &UserConfig{
		Db: DbConfig{
			ExcludeCommands: []string{},
		},
		Gui: GuiConfig{
			ColumnLayout: []Column{
				ColumnExitCode,
				ColumnExecutedAt,
				ColumnCommand,
			},
		},
	}
}
